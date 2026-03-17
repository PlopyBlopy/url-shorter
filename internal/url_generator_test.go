package internal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/tests/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("has suffix", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		ctx := context.Background()

		tg := testdata.NewTestGenerator(0)

		g, err := internal.NewGenerator(tg, ctx)
		require.NoError(err)

		// Act
		url := g.GenerateShortUrl()

		// Assert
		assert.Len(url, len("https://click.ru/XXXXXX"))
	})

	t.Run("no repetitions", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		ctx := context.Background()

		tg := testdata.NewTestGenerator(0)

		g, err := internal.NewGenerator(tg, ctx)
		require.NoError(err)

		// Act
		max := 1069323
		part := 6

		m := make(map[string]int, max)

		repetitionCount := 0

		for i := 0; i < max; i++ {
			url := g.GenerateShortUrl()
			m[url] = m[url] + 1

			if m[url] > 1 {
				repetitionCount++
				fmt.Printf("key:%s, val:%d", url, m[url])
			}

			if part > 1 && i == max/part {
				fmt.Printf("part:6/%d\n", part)
				part--
			}
		}

		// Assert
		assert.Equal(0, repetitionCount)
	})
}
