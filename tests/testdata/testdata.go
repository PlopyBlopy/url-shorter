package testdata

import (
	"context"
	"time"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
)

type TestGenerator struct {
	delta uint64
}

func NewTestGenerator(delta uint64) *TestGenerator {
	return &TestGenerator{
		delta: delta,
	}
}

func (g *TestGenerator) GetCounter(ctx context.Context) (uint64, error) {
	return g.delta, nil
}

// Delta is base counter value generate from internal.generator. In app gets from DB. In testdata from origUrlDelta and shortUrlDelta. CratedAt increment on 1 hour.
func GetUrls(count int, origUrlDelta, shortUrlDelta uint64) ([]domain.Url, error) {
	out := make([]domain.Url, count)

	origGen, err := internal.NewGenerator()
	if err != nil {
		return nil, err
	}

	shortGen, err := internal.NewGenerator()
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		out[i] = domain.Url{
			OrigUrl:   origGen.GenerateShortUrl(uint64(i) + origUrlDelta),
			ShortUrl:  shortGen.GenerateShortUrl(uint64(i) + shortUrlDelta),
			CreatedAt: time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond),
		}
	}

	return out, err
}

func GetUrl(origUrlDelta, shortUrlDelta uint64) (domain.Url, error) {
	url := domain.Url{}

	origGen, err := internal.NewGenerator()
	if err != nil {
		return url, err
	}

	shortGen, err := internal.NewGenerator()
	if err != nil {
		return url, err
	}

	url.OrigUrl = origGen.GenerateShortUrl(origUrlDelta)
	url.ShortUrl = shortGen.GenerateShortUrl(shortUrlDelta)
	url.CreatedAt = time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond)

	return url, nil
}
