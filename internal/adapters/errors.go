package adapters

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func ErrIsUniqueViolates23505(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return false

		}

	}
	return true
}
