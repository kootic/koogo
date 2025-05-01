package dbrepo

import (
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/kootic/koogo/pkg/koohttp"
)

var (
	ErrDBConstraintViolation     = koohttp.NewAPIError(http.StatusConflict, "database_constraint_violation")
	ErrDBTimeout                 = koohttp.NewAPIError(http.StatusRequestTimeout, "database_timeout")
	ErrDBInvalidTransactionState = koohttp.NewAPIError(http.StatusInternalServerError, "database_invalid_transaction_state")
)

func bunErrorHandler(err error) error {
	var pgErr pgdriver.Error

	ok := errors.As(err, &pgErr)
	if !ok {
		return err
	}

	if pgErr.StatementTimeout() {
		return ErrDBTimeout
	}

	if pgErr.IntegrityViolation() {
		return ErrDBConstraintViolation
	}

	// Generally, the above generic error handling is enough as we do not want to
	// expose too much information, but below is an example of how to handle more
	// specific errors.
	switch pgErr.Field('C') {
	case pgerrcode.InvalidTransactionState:
		return ErrDBInvalidTransactionState
	default:
		return err
	}
}
