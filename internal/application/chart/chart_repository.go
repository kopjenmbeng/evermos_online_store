package chart

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

type IChartRepository interface {
	AddChart(ctx context.Context, chart dto.Chart) (status int, err error)
	UpdatChart(ctx context.Context, chart dto.Chart) (status int, err error)
	DeleteChart(ctx context.Context, chart_id string, user_id string) (status int, err error)
	GetChart(ctx context.Context, user_id string) ([]dto.Chart, error)
}

type ChartRepository struct {
	dbr sqlx.QueryerContext
	dbw *sqlx.DB
}

func NewChartRepository(dbr sqlx.QueryerContext, dbw *sqlx.DB) IChartRepository {
	return &ChartRepository{dbr: dbr, dbw: dbw}
}

func (repo *ChartRepository) ValidateStock(ctx context.Context, product_id string, req_qty int) (in_stock bool, err error) {
	query := fmt.Sprintf(`
	select in_stock 
	from public.products
	where product_id=$1 limit 1
	`)
	var stock int = 0
	err = repo.dbw.QueryRowContext(ctx, query, &product_id).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("Barang tidak ditemukan !")
		}
		return false, err
	}
	if stock < req_qty {
		return false, errors.New("Stock tidak cukup !")
	}
	return true, nil
}

func (repo *ChartRepository) GetProductPriceByQty(ctx context.Context, product_id string, qty int) (float64, error) {

	var price float64
	query := fmt.Sprintf(`
	select unit_price 
	from public.products
	where product_id=$1 limit 1
	`)
	err := repo.dbw.QueryRowContext(ctx, query, &product_id).Scan(&price)
	if err != nil {
		return price,err
	}
		
	
	return price * float64(qty), nil
}

func (repo *ChartRepository) AddChart(ctx context.Context, chart dto.Chart) (status int, err error) {
	query := fmt.Sprintf(`
	INSERT INTO public.chart(
		chart_id, qty, total_price, product_id, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6);
	`)
	
	// validate stcok
	inStock, err := repo.ValidateStock(ctx, chart.ProductId,chart.Qty)

	if err != nil {
		return http.StatusBadRequest, err
	}

	if !inStock {
		return http.StatusBadRequest, errors.New("Stock tidak mencukupi")
	}
	// get price base on qty
	price, err := repo.GetProductPriceByQty(ctx, chart.ProductId, chart.Qty)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	chart.TotalPrice = price
	// insert chart
	_, err = repo.dbw.ExecContext(ctx, query,
		&chart.ProductId,
		&chart.Qty,
		&chart.TotalPrice,
		&chart.ProductId,
		&chart.CreatedAt,
		&chart.CreatedBy,
	)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func(repo *ChartRepository)GetProductIdByChartId(ctx context.Context,chart_id string)(string,error){
	var product_id string
	query := fmt.Sprintf(`
	SELECT  product_id
	FROM public.chart where chart_id=$1 limit 1
	`)
	err := repo.dbw.QueryRowContext(ctx, query, &chart_id).Scan(&product_id)
	if err != nil {
		if err==sql.ErrNoRows{
			return product_id,errors.New("Product tidak ditemukan !")
		}
		return product_id,err
	}
		
	
	return product_id, nil
}
func (repo *ChartRepository) UpdatChart(ctx context.Context, chart dto.Chart) (status int, err error) {
	query := fmt.Sprintf(`
	UPDATE public.chart
	SET qty=$1, 
	total_price=$2, 
	updated_at=$3, 
	updated_by=$4
	WHERE created_by=$5 and chart_id=$6 and deleted_by is not null
	
	`)
	// get product id
	product_id,err:=repo.GetProductIdByChartId(ctx,chart.ChartId)
	if err != nil {
		return http.StatusBadRequest, err
	}
	// validate stock
	inStock, err := repo.ValidateStock(ctx, product_id,chart.Qty)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !inStock {
		return http.StatusBadRequest, errors.New("Stock tidak mencukupi")
	}
	// get price base on qty
	price, err := repo.GetProductPriceByQty(ctx, product_id, chart.Qty)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	chart.TotalPrice = price
	_, err = repo.dbw.ExecContext(ctx, query,
		&chart.Qty,
		&chart.TotalPrice,
		&chart.UpdatedAt,
		&chart.UpdatedBy,
		&chart.CreatedBy,
		&chart.ChartId,
	)
	return http.StatusOK, nil
}

func (repo *ChartRepository) DeleteChart(ctx context.Context, chart_id string, user_id string) (status int, err error) {
	query := fmt.Sprintf(`
	UPDATE public.chart
	SET deleted_at=$1, 
	deleted_by=$2
	WHERE created_by=$3 and chart_id=$4
	`)
	now := time.Now()
	_, err = repo.dbw.ExecContext(ctx, query, &now, &user_id, &user_id, &chart_id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusAccepted, nil
}

func (repo *ChartRepository) GetChart(ctx context.Context, user_id string) ([]dto.Chart, error) {
	var result []dto.Chart
	return result, nil
}
