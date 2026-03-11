package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IRepository interface {
	AddCounter(delta int, ctx context.Context) error
	AddCounterTx(delta int, ctx context.Context) error
	IncrementCounter(ctx context.Context) error
	IncrementCounterTx(ctx context.Context) error

	GetCounter(ctx context.Context) (uint64, error)
	GetCounterTx(tx pgx.Tx, ctx context.Context) (uint64, error)

	AddUrl(url domain.Url, ctx context.Context) error
	AddUrls(urls []domain.Url, ctx context.Context) error

	GetUrls(limit int, ctx context.Context) ([]domain.Url, error)
	GetUrlByShortUrl(shortUrl string, ctx context.Context) (domain.Url, error)
	GetUrlByOrigUrl(origUrl string, ctx context.Context) (domain.Url, error)
	GetOrigUrl(shortUrl string, ctx context.Context) (string, error)
	GetShortUrl(origUrl string, ctx context.Context) (string, error)

	RemoveByOrigUrl(origUrl string, ctx context.Context) error
	RemoveByShortUrl(shortUrl string, ctx context.Context)
	RemoveExpired(before time.Time, ctx context.Context) error
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// Увеличивает значение counter на delta.
// Возвращает error DB или pgx.ErrNoRows если строка counter не была затронута.
func (r *Repository) AddCounter(delta int, ctx context.Context) error {
	ct, err := r.pool.Exec(ctx, "UPDATE counter SET counter = counter + $1 WHERE id=0", delta)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Использовать в рамках транзакции.
// Увеличивает значение counter на delta.
// Возвращает error или pgx.ErrNoRows если строка counter не была затронута.
func (r *Repository) AddCounterTx(tx pgx.Tx, delta int, ctx context.Context) error {
	ct, err := tx.Exec(ctx, "UPDATE counter SET counter = counter + $1 WHERE id=0", delta)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Увеличивает значение counter на counter + 1.
// Возвращает ошибку или ошибку pgx.ErrNoRows.
func (r *Repository) IncrementCounter(ctx context.Context) error {
	ct, err := r.pool.Exec(ctx, "UPDATE counter SET counter = counter + 1 WHERE id=0")
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Использовать в рамках транзакции.
// Увеличивает значение counter на counter + 1.
// Возвращает ошибку или ошибку pgx.ErrNoRows.
func (r *Repository) IncrementCounterTx(tx pgx.Tx, ctx context.Context) error {
	ct, err := tx.Exec(ctx, "UPDATE counter SET counter = counter + 1 WHERE id=0")
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Получение значения counter.
// Необходимо проверять на наличие ошибки, так как ошибка pgx.ErrNoRows будет значить, что данные не были прочитаны.
// Возвращает значение counter и error as nil или 0 и error as DB error или pgx.ErrNoRows.
func (r *Repository) GetCounter(ctx context.Context) (uint64, error) {
	var counter uint64

	err := r.pool.QueryRow(ctx, "SELECT counter FROM counter WHERE id = 0").Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// Использовать в рамках транзакции.
// Получение значения counter.
// Необходимо проверять на наличие ошибки, так как ошибка pgx.ErrNoRows будет значить, что данные не были прочитаны.
// Возвращает значение counter и error as nil или 0 и error as DB error или pgx.ErrNoRows.
func (r *Repository) GetCounterTx(tx pgx.Tx, ctx context.Context) (uint64, error) {
	var counter uint64

	err := tx.QueryRow(ctx, "SELECT counter FROM counter WHERE id = 0").Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// Реализованна как транзакция.
// Добавление структуры domain.Url в базу данных.
// Автоматически увеличивает counter.
func (r *Repository) AddUrl(url domain.Url, ctx context.Context) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, "INSERT INTO urls (orig_url, short_url, created_at) VALUES ($1, $2, $3)", url.OrigUrl, url.ShortUrl, url.CreatedAt)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	err = r.IncrementCounterTx(tx, ctx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Реализованна как транзакция.
// Добавление []domain.Url в базу данных.
// Автоматически увеличивает counter на delta=len(urls).
// Может вернуть ошибку или ошибку pgx.ErrNoRows
func (r *Repository) AddUrls(urls []domain.Url, ctx context.Context) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(urls) < 50 {
		err = r.addUrlsPackage(tx, urls, ctx)
		if err != nil {
			return err
		}
	} else {
		err = r.addUrlsCopyFrom(tx, urls, ctx)
		if err != nil {
			return err
		}
	}

	err = r.AddCounterTx(tx, len(urls), ctx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Вспомогательная функция для пакетного добавления []domain.Url
// Может вернуть ошибку или ошибку pgx.ErrNoRows
func (r *Repository) addUrlsPackage(tx pgx.Tx, urls []domain.Url, ctx context.Context) error {
	var sb strings.Builder
	args := make([]any, len(urls)*3)

	sb.WriteString("INSERT INTO urls (orig_url, short_url, created_at) VALUES ")

	numCols := 3
	for i, u := range urls {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j := 0; j < numCols; j++ {
			if j > 0 {
				sb.WriteString(", ")
			}
			fmt.Fprintf(&sb, "$%d", i*numCols+j+1)
		}
		sb.WriteString(")")

		base := i * 3
		args[base] = u.OrigUrl
		args[base+1] = u.ShortUrl
		args[base+2] = u.CreatedAt
	}

	ct, err := tx.Exec(ctx, sb.String(), args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Вспомогательная функция для []domain.Url через протокол COPY в postgres
// Может вернуть ошибку или ошибку pgx.ErrNoRows
func (r *Repository) addUrlsCopyFrom(tx pgx.Tx, urls []domain.Url, ctx context.Context) error {
	copyCount, err := tx.CopyFrom(ctx, pgx.Identifier{"urls"}, []string{"orig_url", "short_url", "created_at"}, pgx.CopyFromSlice(len(urls), func(i int) ([]any, error) {
		u := urls[i]
		return []any{
			u.OrigUrl,
			u.ShortUrl,
			u.CreatedAt,
		}, nil

	}))
	if err != nil {
		return err
	}
	if copyCount != int64(len(urls)) {
		return pgx.ErrNoRows
	}

	return nil
}

// Получение domain.Url где их кол-во = limit.
// Может вернуть []domain.Url с элементами или empty []domain.Url.
func (r *Repository) GetUrls(limit int, ctx context.Context) ([]domain.Url, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM urls LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	urls, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Url])
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// Получение структуры domain.Url по shortUrl
// Может вернуть domain.Url, error=nil или default domain.Url, error
func (r *Repository) GetUrlByShortUrl(shortUrl string, ctx context.Context) (domain.Url, error) {
	var url domain.Url

	err := r.pool.QueryRow(ctx, "SELECT * FROM urls WHERE short_url = $1", shortUrl).Scan(&url.OrigUrl, &url.ShortUrl, &url.CreatedAt)
	if err != nil {
		return url, err
	}

	return url, nil
}

// Получение структуры domain.Url по origUrl
// Может вернуть domain.Url, error=nil или default domain.Url, error
func (r *Repository) GetUrlByOrigUrl(origUrl string, ctx context.Context) (domain.Url, error) {
	var url domain.Url

	err := r.pool.QueryRow(ctx, "SELECT * FROM urls WHERE orig_url = $1", origUrl).Scan(&url.OrigUrl, &url.ShortUrl, &url.CreatedAt)
	if err != nil {
		return url, err
	}

	return url, nil
}

// Получение string = origUrl по shortUrl
// Может вернуть origUrl, error=nil или default origUrl, error
func (r *Repository) GetOrigUrl(shortUrl string, ctx context.Context) (string, error) {
	var origUrl string

	err := r.pool.QueryRow(ctx, "SELECT orig_url FROM urls WHERE short_url = $1", shortUrl).Scan(&origUrl)
	if err != nil {
		return "", err
	}
	return origUrl, nil
}

// Получение string = shortUrl по origUrl
// Может вернуть shortUrl, error=nil или default shortUrl, error
func (r *Repository) GetShortUrl(origUrl string, ctx context.Context) (string, error) {

	var shortUrl string
	err := r.pool.QueryRow(ctx, "SELECT short_url FROM urls WHERE orig_url = $1", origUrl).Scan(&shortUrl)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

// Удаление row по origUrl
// Может вернуть nil или error или error = pgx.ErrNoRows
func (r *Repository) RemoveByOrigUrl(origUrl string, ctx context.Context) error {

	ct, err := r.pool.Exec(ctx, "DELETE FROM urls WHERE orig_url = $1", origUrl)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Удаление row по shortUrl
// Может вернуть nil или error или error = pgx.ErrNoRows
func (r *Repository) RemoveByShortUrl(shortUrl string, ctx context.Context) error {
	ct, err := r.pool.Exec(ctx, "DELETE FROM urls WHERE short_url = $1", shortUrl)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Удаление row по истечению времени
// Может вернуть nil или error или error = pgx.ErrNoRows
func (r *Repository) RemoveExpired(before time.Time, ctx context.Context) error {
	ct, err := r.pool.Exec(ctx, "DELETE FROM urls WHERE created_at < $1", before)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
