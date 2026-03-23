package internal

// The maximum number of unique short links.
// Collisions are possible if this value is exceeded.
const maxUnique = 106_932_384

type Generator struct {
}

func NewGenerator() (*Generator, error) {
	return &Generator{}, nil
}

// TODO: в случае если counter не был увеличен на стороне БД, то при следующем запуске приложения создастся короткая ссылка - что уже есть в БД. Происходит из за того что есть внутренний counter и counter на стороне БД - несвязанные во время жизни приложения
func (g *Generator) GenerateShortUrl(counter uint64) string {
	buf := []byte{'h', 't', 't', 'p', 's', ':', '/', '/', 'c', 'l', 'i', 'c', 'k', '.', 'r', 'u', '/', ' ', ' ', ' ', ' ', ' ', ' '}

	if counter >= maxUnique {
		panic("The indicator value limit has been reached")
	}

	n := counter

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
