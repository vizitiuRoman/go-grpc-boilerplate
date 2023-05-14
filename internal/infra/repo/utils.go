package repo

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain"
)

func wrapErrNoRows(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func isUniqueViolation(err error) bool {
	pgErr, ok := err.(*pq.Error)
	return ok && pgErr.Code == "23505"
}

func wrapUniqueViolation(err error) error {
	if isUniqueViolation(err) {
		return domain.ErrAlreadyExists
	}
	return err
}
