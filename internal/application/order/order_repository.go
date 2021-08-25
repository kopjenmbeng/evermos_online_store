package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kopjenmbeng/evermos_online_store/internal/dto"
)

type IOrderRepository interface {
	Create(ctx context.Context, odrs []dto.Order) (status int, err error)
}

type OrderRepository struct {
	dbr sqlx.QueryerContext
	dbw *sqlx.DB
}

func NewOrderRepository(dbr sqlx.QueryerContext, dbw *sqlx.DB) IOrderRepository {
	return &OrderRepository{dbr: dbr, dbw: dbw}
}

func (repo *OrderRepository) Create(ctx context.Context, odrs []dto.Order) (status int, Err error) {

	query := fmt.Sprintf(`
	INSERT INTO public.orders(
		order_id, status, created_at, created_by)
		VALUES ($1, $2, $3, $4);
	`)
	tx, err := repo.dbw.BeginTx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, odr := range odrs {

		// get product id
		product_id, err := repo.GetProductIdByChart(ctx, odr.OrderId, tx)
		if err != nil {
			return http.StatusBadRequest, err
		}

		// check stock
		inStock, remain, err := repo.CheckStock(ctx, odr.OrderId, tx)
		if err != nil {
			return http.StatusBadRequest, err
		}
		if !inStock {
			return http.StatusBadRequest, errors.New("Stock tidak cukup !")
		}
		// create order
		_, err = tx.ExecContext(ctx, query,
			&odr.OrderId,
			&odr.Status,
			&odr.CreatedAt,
			&odr.CreatedBy,
		)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// update stock
		status, err = repo.UpdateStock(ctx, product_id, remain, tx)
		if err != nil {
			tx.Rollback()
			return status,err
		}
		// update order status
		status, err = repo.UpdateChartStatus(ctx, odr.OrderId, odr.CreatedBy, tx)
		if err != nil {
			tx.Rollback()
			return status,err
		}
	}

	// commit all transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	return
}

func (repo *OrderRepository) GetProductIdByChart(ctx context.Context, chart_id string, tx *sql.Tx) (string, error) {
	var product_id string = ""
	query := fmt.Sprintf(`
	SELECT product_id
	FROM public.chart where chart_id=$1 limit 1
	`)
	err := tx.QueryRowContext(ctx, query, &chart_id).Scan(&product_id)
	if err != nil {
		return "", err
	}
	return product_id, nil
}
func (repo *OrderRepository) CheckStock(ctx context.Context, chart_id string, tx *sql.Tx) (bool, int, error) {
	var remain int
	query := fmt.Sprintf(`
	select (in_stock - ch.qty) as sisa
	from public.products p inner join 
	chart ch on ch.product_id=p.product_id
	where ch.chart_id=$1 limit 1
	`)

	err := tx.QueryRowContext(ctx, query, &chart_id).Scan(&remain)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, errors.New("Stock tidak.")
		}
		return false, 0, err
	}
	if remain < 0 {
		return false, 0, errors.New("Stok tidak cukup !.")
	}
	return true, remain, nil
}
func (repo *OrderRepository) UpdateStock(ctx context.Context, product_id string, remain int, tx *sql.Tx) (status int, err error) {
	query := fmt.Sprintf(`
	UPDATE public.products 
	SET in_stock=$1
	WHERE product_id=$2
	`)
	_, err = tx.ExecContext(ctx, query,
		&remain,
		&product_id,
	)
	if err != nil {
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("error at step update stock %s", err.Error()))
	}
	return http.StatusCreated, nil
}

func (repo *OrderRepository) UpdateChartStatus(ctx context.Context, chart_id string, user_id string, tx *sql.Tx) (status int, err error) {
	query := fmt.Sprintf(`
	UPDATE public.chart
	SET  updated_at=$1, 
	updated_by=$2,
	is_ordered='true'
	WHERE chart_id=$3 and created_by=$4
	`)
	now := time.Now()
	_, err = tx.ExecContext(ctx, query,
		&now,
		&user_id,
		&chart_id,
		&user_id,
	)
	if err != nil {
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("error at step update chart status %s", err.Error()))
	}
	return http.StatusCreated, nil
}
