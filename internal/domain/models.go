package domain

import "time"

type Url struct {
	OrigUrl   string    `json:"origUrl" db:"orig_url"`
	ShortUrl  string    `json:"shortUrl" db:"short_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
