package order

import (
	"context"
	"net/http"
	"time"

	"github.com/kopjenmbeng/evermos_online_store/internal/dto"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
)

type IOrderUseCase interface {
	Create(ctx context.Context, req CreateOrderRequest) (status int, err error)
}

type OrderUseCase struct {
	repository IOrderRepository
	r          *http.Request
}

func NewOrderUserCase(repo IOrderRepository,r *http.Request)IOrderUseCase{
	return &OrderUseCase{repository: repo,r: r}
}

func(use_case *OrderUseCase)Create(ctx context.Context, req CreateOrderRequest) (status int, err error){
	claim:=jwe_auth.GetClaims(use_case.r)
	for _,i:= range req.Charts{
		order:=dto.Order{OrderId: i,Status: "Menunggu Konfirmasi",CreatedAt: time.Now(),CreatedBy: claim.Public.Subject}
		status,err=use_case.repository.Create(ctx,order)
		if err!=nil{
			return
		}
	}
	return
	
}
