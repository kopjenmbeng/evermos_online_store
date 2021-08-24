package chart

import (
	"net/http"

	"github.com/kopjenmbeng/evermos_online_store/internal/utility/validator"
)

type UpdateChartRequest struct {
	ChartId string `json:"chart_id"`
	Qty       int    `json:"qty"`
}

func (req *UpdateChartRequest) Bind(r *http.Request) error {
	if err := req.Validate(r); err != nil {
		return err
	}
	return nil
}

func (req *UpdateChartRequest) Validate(r *http.Request) error {
	if err := validator.ValidateEmpty("chart_id", req.ChartId); err != nil {
		return err
	}
	if err := validator.ValidateEmpty("qty", req.Qty); err != nil {
		return err
	}
	return nil
}
