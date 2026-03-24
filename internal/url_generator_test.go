package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("has suffix", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		g, err := NewGenerator()
		require.NoError(err)

		// Act
		url := g.GenerateShortUrl(uint64(0))

		// Assert
		assert.Len(url, len("https://click.ru/XXXXXX"))
	})

	t.Run("the maximum unique value limit has been reached", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		g, err := NewGenerator()
		require.NoError(err)

		var max uint64 = maxUnique + 1

		// Assert
		assert.Panics(func() {
			_ = g.GenerateShortUrl(max)
		})
	})
}

// The test may end earlier due to the end of the timeout for the test.
// It needs to run in debug or with an extended timeout.
func TestNoRepetitionToMaxUniqueValue(t *testing.T) {
	// Arrange
	require := require.New(t)
	assert := assert.New(t)

	g, err := NewGenerator()
	require.NoError(err)

	var max uint64 = maxUnique

	part := uint64(6)

	m := make(map[string]int, max)
	repetitionCount := 0
	var val uint64

	// Act
	for val = 0; val < max; val++ {
		url := g.GenerateShortUrl(val)
		m[url] = m[url] + 1

		if m[url] > 1 {
			repetitionCount++
			fmt.Printf("key:%s, val:%d", url, m[url])
		}

		if part > 0 && val == max/part {
			fmt.Printf("part:6/%d\n", part)
			part--
		}
	}

	// Assert
	assert.Equal(0, repetitionCount)
}
