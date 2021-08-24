package authentication

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kopjenmbeng/evermos_online_store/internal/dto"
)

type IAuthenticationRepository interface {
	GetUser(ctx context.Context, phone_number string) (data *dto.Customer, status int, Err error)
	Register(ctx context.Context, customer dto.Customer) (status int, err error)
}

type AuthenticationRepository struct {
	dbr sqlx.QueryerContext
	dbw *sqlx.DB
}

func NewAuthenticationRepository(dbr sqlx.QueryerContext, dbw *sqlx.DB) IAuthenticationRepository {
	return &AuthenticationRepository{dbr: dbr, dbw: dbw}
}

func (repo *AuthenticationRepository) IsDuplicatePhoneNumber(ctx context.Context, phone_number string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT count(*) as total
	FROM public.customers where phone_number=$1 limit 1
	`)

	var total int = 0
	err := repo.dbw.QueryRowContext(ctx, query, &phone_number).Scan(&total)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}

	return false, nil

}
func (repo *AuthenticationRepository) Register(ctx context.Context, customer dto.Customer) (status int, err error) {

	query := fmt.Sprintf(`
	INSERT INTO public.customers(
		customer_id, phone_number, full_name, salt, password, iteration, security_length)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`)

	isDuplicate,err:=repo.IsDuplicatePhoneNumber(ctx,customer.PhoneNumber)
	if err!=nil{
		return http.StatusInternalServerError,err
	}
	if isDuplicate{
		return http.StatusBadRequest,errors.New("Nomor telp sudah terdaftar !.")
	}
	// usr:=dto.Customer{CustomerId: uuid.New().String(),PhoneNumber: phone_number,FullName: full_name,Salt: salt,Password: password,Iteration: }
	_, err = repo.dbw.ExecContext(ctx, query,
		&customer.CustomerId,
		&customer.PhoneNumber,
		&customer.FullName,
		&customer.Salt,
		&customer.Password,
		&customer.Iteration,
		&customer.SecurityLength,
	)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}
	status = http.StatusCreated
	return
}
func (repo *AuthenticationRepository) GetUser(ctx context.Context, phone_number string) (data *dto.Customer, status int, Err error) {
	query := fmt.Sprintf(`
	SELECT customer_id, phone_number, full_name, salt, password, iteration, security_length
	FROM public.customers where phone_number=$1 limit 1
	`)

	var cus dto.Customer
	err := repo.dbr.QueryRowxContext(ctx, query, &phone_number).Scan(&cus.CustomerId, &phone_number, &cus.FullName, &cus.Salt, &cus.Password, &cus.Iteration, &cus.SecurityLength)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusUnauthorized, errors.New("Nomor telp belum terdaftar !")
		}
		return nil, http.StatusInternalServerError, err
	}
	return &cus, http.StatusOK, nil
}
