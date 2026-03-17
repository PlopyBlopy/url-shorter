package testdata

import (
	"context"
	"time"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
)

type testGenerator struct {
	delta uint64
}

func (g *testGenerator) GetCounter(ctx context.Context) (uint64, error) {
	return g.delta, nil
}

// Delta is base counter value generate from internal.generator. In app gets from DB. In testdata from origUrlDelta and shortUrlDelta. CratedAt increment on 1 hour.
func GetUrls(count int, origUrlDelta, shortUrlDelta uint64) ([]domain.Url, error) {
	out := make([]domain.Url, count)
	g1 := &testGenerator{delta: origUrlDelta}
	g2 := &testGenerator{delta: shortUrlDelta}

	ctx := context.Background()

	origGen, err := internal.NewGenerator(g1, ctx)
	if err != nil {
		return nil, err
	}

	shortGen, err := internal.NewGenerator(g2, ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		out[i] = domain.Url{
			OrigUrl:   origGen.GenerateShortUrl(),
			ShortUrl:  shortGen.GenerateShortUrl(),
			CreatedAt: time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond),
		}
	}

	return out, err
}

func GetUrl(origUrlDelta, shortUrlDelta uint64) (domain.Url, error) {
	g1 := &testGenerator{delta: origUrlDelta}
	g2 := &testGenerator{delta: shortUrlDelta}

	ctx := context.Background()

	url := domain.Url{}

	origGen, err := internal.NewGenerator(g1, ctx)
	if err != nil {
		return url, err
	}

	shortGen, err := internal.NewGenerator(g2, ctx)
	if err != nil {
		return url, err
	}

	url.OrigUrl = origGen.GenerateShortUrl()
	url.ShortUrl = shortGen.GenerateShortUrl()
	url.CreatedAt = time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond)

	return url, nil
}
