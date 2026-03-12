///go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/PlopyBlopy/url-shorter/tests/testdata"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	Все тесты делались по системе AAA (Arrange, Act, Assert)
*/

func TestWithPostgresTestcontainerAndWithRollbackShapshotOnTest(t *testing.T) {
	// create postgres container
	pgContainerCtx := context.Background()
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(pgContainerCtx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	// create TestSuite
	testSuite, err := testdata.NewTestSuite(pgContainer, pgContainerCtx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("IncrementCounter", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err := testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		//Act
		expCounter, err := rep.GetCounter(ctx)
		assert.NoError(err)

		err = rep.IncrementCounter(ctx)
		assert.NoError(err)

		actCounter, err := rep.GetCounter(ctx)
		assert.NoError(err)

		//Assert
		assert.Equalf(expCounter, actCounter-1, "exp=%d, act=%d", expCounter, actCounter-1)
		assert.NotEqualf(expCounter, actCounter, "exp=%d, act=%d", expCounter, actCounter)
	})

	t.Run("IncrementCounterTx", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err := testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		//Act
		tx, err := testSuite.Db.Begin(ctx)
		assert.NoError(err)
		defer tx.Rollback(ctx)

		expCounter, err := rep.GetCounterTx(tx, ctx)
		assert.NoError(err)

		err = rep.IncrementCounterTx(tx, ctx)
		assert.NoError(err)

		actCounter, err := rep.GetCounterTx(tx, ctx)
		assert.NoError(err)

		err = tx.Commit(ctx)
		assert.NoError(err)

		//Assert
		assert.Equalf(expCounter, actCounter-1, "exp=%d, act=%d", expCounter, actCounter-1)
		assert.NotEqualf(expCounter, actCounter, "exp=%d, act=%d", expCounter, actCounter)
	})

	t.Run("AddCounter", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err := testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		//Act
		expCounter, err := rep.GetCounter(ctx)
		assert.NoError(err)

		err = rep.AddCounter(1, ctx)
		assert.NoError(err)

		actCounter, err := rep.GetCounter(ctx)
		assert.NoError(err)

		//Assert
		assert.Equalf(expCounter, actCounter-1, "exp=%d, act=%d", expCounter, actCounter-1)
		assert.NotEqualf(expCounter, actCounter, "exp=%d, act=%d", expCounter, actCounter)
	})

	t.Run("AddCounterTx", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err := testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		//Act
		tx, err := testSuite.Db.Begin(ctx)
		assert.NoError(err)
		defer tx.Rollback(ctx)

		expCounter, err := rep.GetCounterTx(tx, ctx)
		assert.NoError(err)

		err = rep.AddCounterTx(tx, 1, ctx)
		assert.NoError(err)

		actCounter, err := rep.GetCounterTx(tx, ctx)
		assert.NoError(err)

		err = tx.Commit(ctx)
		assert.NoError(err)

		//Assert
		assert.Equalf(expCounter, actCounter-1, "exp=%d, act=%d", expCounter, actCounter-1)
		assert.NotEqualf(expCounter, actCounter, "exp=%d, act=%d", expCounter, actCounter)
	})

	t.Run("GetCounter", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err := testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		var expCounter uint64

		//Act
		actCounter, err := rep.GetCounter(ctx)
		assert.NoError(err)

		//Assert
		assert.Equalf(expCounter, actCounter, "exp=%d, act=%d", expCounter, actCounter)
	})

	t.Run("AddUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(pgContainerCtx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		origUrl := "https://testdomain/testroute"
		shortUrl := gen.GenerateShortUrl()
		now := time.Now().UTC().Truncate(time.Microsecond)

		expUrl := domain.Url{
			OrigUrl:   origUrl,
			ShortUrl:  shortUrl,
			CreatedAt: now,
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		actUrl, err := rep.GetUrlByShortUrl(shortUrl, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrl, actUrl)
	})

	t.Run("addUrlsPackage", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		urlsCount := 49

		expUrls, err := testdata.GetUrls(urlsCount, 0, uint64(urlsCount))
		assert.NoError(err)

		//Act
		err = rep.AddUrls(expUrls, ctx)
		assert.NoError(err)

		actUrls, err := rep.GetUrls(urlsCount, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrls, actUrls)
		assert.NotEqual(expUrls[1:], actUrls)
	})

	t.Run("addUrlsCopyFrom", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		urlsCount := 50

		expUrls, err := testdata.GetUrls(urlsCount, 0, uint64(urlsCount))
		assert.NoError(err)

		//Act
		err = rep.AddUrls(expUrls, ctx)
		assert.NoError(err)

		actUrls, err := rep.GetUrls(urlsCount, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrls, actUrls)
		assert.NotEqual(expUrls[1:], actUrls)
	})

	t.Run("GetUrls", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		tests := []struct {
			limit     int
			expOutput int
		}{
			{0, 0},
			{1, 1},
			{10, 10},
		}

		expUrls, err := testdata.GetUrls(10, 0, 10)
		assert.NoError(err)

		//Act
		err = rep.AddUrls(expUrls, ctx)
		assert.NoError(err)

		//Assert
		for _, test := range tests {
			actUrls, err := rep.GetUrls(test.limit, ctx)
			assert.NoError(err)

			assert.Equal(test.expOutput, len(actUrls))
			assert.Equal(expUrls[:test.expOutput], actUrls)
		}

	})
	t.Run("GetUrlWithShortUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		actUrl, err := rep.GetUrlByShortUrl(expUrl.ShortUrl, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrl, actUrl)
	})
	t.Run("GetUrlWithOrigUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		actUrl, err := rep.GetUrlByOrigUrl(expUrl.OrigUrl, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrl, actUrl)
	})
	t.Run("GetShortUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		actUrl, err := rep.GetShortUrl(expUrl.OrigUrl, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrl.ShortUrl, actUrl)
	})
	t.Run("GetOrigUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		actUrl, err := rep.GetOrigUrl(expUrl.ShortUrl, ctx)
		assert.NoError(err)

		//Assert
		assert.Equal(expUrl.OrigUrl, actUrl)
	})

	t.Run("RemoveByOrigUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		err = rep.RemoveByOrigUrl(expUrl.OrigUrl, ctx)
		assert.NoError(err)

		actUrl, actErr := rep.GetUrlByOrigUrl(expUrl.OrigUrl, ctx)

		//Assert
		assert.Condition(func() bool {
			defaultUrl := domain.Url{}
			return actUrl == defaultUrl
		})
		assert.Error(actErr)
		assert.ErrorIs(actErr, pgx.ErrNoRows)
	})
	t.Run("RemoveByShortUrl", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		expUrl := domain.Url{
			OrigUrl:   "https://testdomain/testroute",
			ShortUrl:  gen.GenerateShortUrl(),
			CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		}

		//Act
		err = rep.AddUrl(expUrl, ctx)
		assert.NoError(err)

		err = rep.RemoveByShortUrl(expUrl.ShortUrl, ctx)
		assert.NoError(err)

		actUrl, actErr := rep.GetUrlByShortUrl(expUrl.ShortUrl, ctx)

		//Assert
		assert.Condition(func() bool {
			defaultUrl := domain.Url{}
			return actUrl == defaultUrl
		})
		assert.Error(actErr)
		assert.ErrorIs(actErr, pgx.ErrNoRows)
	})
	t.Run("RemoveExpired", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		err = testSuite.SetupTest(ctx)
		require.NoError(t, err)

		assert := assert.New(t)

		rep := adapters.NewRepository(testSuite.Db)

		gen, err := internal.NewGenerator(rep, ctx)
		assert.NoError(err)

		tests := []struct {
			urls []domain.Url
		}{
			{[]domain.Url{
				{
					OrigUrl:   "https://testdomain/testroute_0",
					ShortUrl:  gen.GenerateShortUrl(),
					CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
				},
			}},
			{[]domain.Url{
				{
					OrigUrl:   "https://testdomain/testroute_0",
					ShortUrl:  gen.GenerateShortUrl(),
					CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
				},
			}},
			{[]domain.Url{
				{
					OrigUrl:   "https://testdomain/testroute_0",
					ShortUrl:  gen.GenerateShortUrl(),
					CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
				},
				{
					OrigUrl:   "https://testdomain/testroute_1",
					ShortUrl:  gen.GenerateShortUrl(),
					CreatedAt: time.Now().UTC().Truncate(time.Microsecond).Add(time.Hour * 2),
				},
			}},
		}

		//Act
		for _, test := range tests {
			for _, url := range test.urls {
				err = rep.AddUrl(url, ctx)
				assert.NoError(err)
			}

			err = rep.RemoveExpired(test.urls[0].CreatedAt.Add(time.Minute), ctx)
			assert.NoError(err)

			actUrl, actErr := rep.GetUrlByShortUrl(test.urls[0].ShortUrl, ctx)

			//Assert
			assert.Condition(func() bool {
				defaultUrl := domain.Url{}
				return actUrl == defaultUrl
			})
			assert.Error(actErr)
			assert.ErrorIs(actErr, pgx.ErrNoRows)
		}

	})

}

func TestWithPostgresTestcontainerAndWithoutRollbackShapshotOnTest(t *testing.T) {
	// create postgres container
	pgContainerCtx := context.Background()
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(pgContainerCtx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	// create TestSuite
	testSuite, err := testdata.NewTestSuite(pgContainer, pgContainerCtx)
	if err != nil {
		t.Fatal(err)
	}

	// Test here
	_ = testSuite // placeholder - replace him
}
