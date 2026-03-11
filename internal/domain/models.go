package domain

import "time"

type Url struct {
	OrigUrl   string    `json:"origUrl" db:"orig_url"`
	ShortUrl  string    `json:"shortUrl" db:"short_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type OrigUrl struct {
	OrigUrl string `json:"url" binding:"required,url"`
}

type ShortUrl struct {
	ShortUrl string `json:"url" binding:"required,url"`
}

func NewUrl(origUrl, shortUrl string) Url {
	return Url{
		OrigUrl:   origUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
}
