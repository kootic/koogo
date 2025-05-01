package postgres

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/kootic/koogo/pkg/koohttp"
)

var (
	ErrNotFound                = koohttp.NewAPIError(http.StatusNotFound, "database_record_not_found")
	ErrConstraintViolation     = koohttp.NewAPIError(http.StatusConflict, "database_constraint_violation")
	ErrTimeout                 = koohttp.NewAPIError(http.StatusRequestTimeout, "database_timeout")
	ErrInvalidTransactionState = koohttp.NewAPIError(http.StatusInternalServerError, "database_invalid_transaction_state")
)

func handleError(err error) error {
	// Handle sql.ErrNoRows - record not found
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr pgdriver.Error

	ok := errors.As(err, &pgErr)
	if !ok {
		return err
	}

	if pgErr.StatementTimeout() {
		return ErrTimeout
	}

	if pgErr.IntegrityViolation() {
		return ErrConstraintViolation
	}

	// Generally, the above generic error handling is enough as we do not want to
	// expose too much information, but below is an example of how to handle more
	// specific errors.
	switch pgErr.Field('C') {
	case pgerrcode.InvalidTransactionState:
		return ErrInvalidTransactionState
	default:
		return err
	}
}
