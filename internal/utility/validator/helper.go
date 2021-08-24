package validator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

func ValidateEmpty(fieldName string, param interface{}) error {
	return validation.Errors{
		fieldName: validation.Validate(&param, validation.Required),
	}.Filter()
}

func ValidateExistInDatabase(fieldName string, param interface{}, tableName string, dbr sqlx.QueryerContext, ctx context.Context) error {
	var (
		total sql.NullInt32
	)

	selectQuery := fmt.Sprintf(`
	SELECT COUNT(%s) FROM %s as tbl WHERE %s=$1
	`, fieldName, tableName, fieldName)

	if err := dbr.QueryRowxContext(ctx, selectQuery, &param).Scan(&total); err != nil {
		return errors.New(fmt.Sprintf("value %s Does'nt Exist !", param))
	}

	if total.Int32 == 0 {
		return errors.New(fmt.Sprintf("value %s Does'nt Exist !", param))
	}

	return nil
}
