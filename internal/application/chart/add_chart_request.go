package chart

import (
	"net/http"

	"github.com/kopjenmbeng/evermos_online_store/internal/utility/validator"
)

type AddChartRequest struct {
	ProductId string `json:"product_id"`
	Qty       int    `json:"qty"`
}

func (req *AddChartRequest) Bind(r *http.Request) error {
	if err := req.Validate(r); err != nil {
		return err
	}
	return nil
}

func (req *AddChartRequest) Validate(r *http.Request) error {
	if err := validator.ValidateEmpty("product_id", req.ProductId); err != nil {
		return err
	}
	if err := validator.ValidateEmpty("qty", req.Qty); err != nil {
		return err
	}
	return nil
}
