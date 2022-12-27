package kun

import (
	"database/sql"

	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Error error

var (
	// ErrNotFound error
	ErrNotFound Error = errors.New("not found")

	// ErrAlreadyExists error
	ErrAlreadyExists Error = errors.New("already exists")
)

// Handle error
func HandleError(err error) error {
	if err != nil {
		if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() {
			switch pgErr.Field('C') {
			case pgerrcode.NoDataFound:
				return ErrNotFound
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
	}
	return errors.Wrap(err, "not handled pg error")
}
