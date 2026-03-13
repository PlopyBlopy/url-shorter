package domain

import (
	"context"
)

type ShortURLGenerator interface {
	GenerateShortUrl() string
}

type UrlAddGetter interface {
	UrlAdder
	UrlShortGetter
}

type UrlAdder interface {
	AddUrl(url Url, ctx context.Context) error
}

type UrlShortGetter interface {
	GetShortUrl(origUrl string, ctx context.Context) (string, error)
}

type UrlOrigGetter interface {
	GetOrigUrl(shortUrl string, ctx context.Context) (string, error)
}

type UrlsGetter interface {
	GetUrls(limit int, ctx context.Context) ([]Url, error)
}
