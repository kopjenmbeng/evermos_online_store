package chart

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kopjenmbeng/evermos_online_store/internal/dto"
	// "github.com/kopjenmbeng/evermos_online_store/internal/middleware"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
)

type IChartUseCase interface {
	AddToChart(ctx context.Context, req AddChartRequest) (status int, err error)
	UpdateChart(ctx context.Context, req UpdateChartRequest) (status int, err error)
}

type ChartUseCase struct {
	repository IChartRepository
	r          *http.Request
}

func NewChartUseCase(repo IChartRepository, r *http.Request) IChartUseCase {
	return &ChartUseCase{repository: repo, r: r}
}

func (use_case *ChartUseCase) AddToChart(ctx context.Context, req AddChartRequest) (status int, err error) {

	claim:=jwe_auth.GetClaims(use_case.r)
	chart:=dto.Chart{ChartId: uuid.New().String(),ProductId: req.ProductId,Qty: req.Qty,CreatedAt: time.Now(),CreatedBy: claim.Public.Subject}
	status,err=use_case.repository.AddChart(ctx,chart)
	return
}

func(use_case *ChartUseCase)UpdateChart(ctx context.Context, req UpdateChartRequest) (status int, err error){
	claim:=jwe_auth.GetClaims(use_case.r)
	chart:=dto.Chart{ChartId: req.ChartId,Qty: req.Qty,UpdatedAt: time.Now(),UpdatedBy: claim.Public.Subject,CreatedBy: claim.Public.Subject}
	status,err=use_case.repository.UpdatChart(ctx,chart)
	return
}