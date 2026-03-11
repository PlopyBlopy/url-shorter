package addUrl

import (
	"context"
	"errors"
	"fmt"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/jackc/pgx/v5"
)

func AddUrlUsecase(g *internal.Generator, rep *adapters.Repository) func(string, context.Context) (string, error) {
	return func(origUrl string, ctx context.Context) (string, error) {
		shortUrl, err := rep.GetShortUrl(origUrl, ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				shortUrl = g.GenerateShortUrl()

				url := domain.NewUrl(origUrl, shortUrl)

				err = rep.AddUrl(url, ctx)
				if err != nil {
					return "", fmt.Errorf("Failed add url: %w", err)
				}
			} else {
				return "", err
			}
		}

		return shortUrl, nil
	}
}
