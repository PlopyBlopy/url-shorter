package internal

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
)

type Generator struct {
	counter uint64
	rep     domain.CounterGetter
}

func NewGenerator(rep domain.CounterGetter, ctx context.Context) (*Generator, error) {
	counter, err := rep.GetCounter(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Generator: %w", err)
	}

	return &Generator{
		counter: counter,
		rep:     rep,
	}, nil
}

func (g *Generator) GenerateShortUrl() string {
	buf := []byte{'h', 't', 't', 'p', 's', ':', '/', '/', 'c', 'l', 'i', 'c', 'k', '.', 'r', 'u', '/', ' ', ' ', ' ', ' ', ' ', ' '}

	n := atomic.AddUint64(&g.counter, 1) - 1 // starts from 0 if first, or couter num

	if n == math.MaxUint64 {
		panic("The indicator value limit has been reached")
	}

	const (
		base26 = 26
		base9  = 9
	)

	buf[len(buf)-6] = byte('A' + n%base26)
	n /= base26

	buf[len(buf)-5] = byte('A' + n%base26)
	n /= base26

	buf[len(buf)-4] = byte('1' + n%base9)
	n /= base9

	buf[len(buf)-3] = byte('a' + n%base26)
	n /= base26

	buf[len(buf)-2] = byte('A' + n%base26)
	n /= base26

	buf[len(buf)-1] = byte('A' + n%base26)
	n /= base26

	return string(buf)
}
