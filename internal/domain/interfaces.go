package domain

import (
	"context"
)

type ShortURLGenerator interface {
	GenerateShortUrl(uint64) string
}

type UrlAddGetter interface {
	UrlAdder
	UrlShortGetter
}

type UrlAdder interface {
	AddUrl(url Url, ctx context.Context) error
}
type UrlAnyGetter interface {
	GetUrl(anyurl string, ctx context.Context) (Url, error)
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

type CounterGetter interface {
	GetCounter(ctx context.Context) (uint64, error)
}
