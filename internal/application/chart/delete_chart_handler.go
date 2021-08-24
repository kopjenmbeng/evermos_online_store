package chart

import (
	"net/http"

	"github.com/RoseRocket/xerrs"
	"github.com/go-chi/chi"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware"
	"github.com/kopjenmbeng/evermos_online_store/internal/utility/respond"
	"github.com/kopjenmbeng/evermos_online_store/internal/utility/validator"
)

func DeleteChartHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// err error
		rc  = r.Context()
		err error
	)

	ChartId := chi.URLParam(r, "chart_id")
	if err := validator.ValidateEmpty("chart_id", ChartId); err != nil {
		respond.Nay(w, r, http.StatusBadRequest, err)
		return
	}
	useCase := UseCaseFromContext(rc)
	code, err := useCase.DeleteChart(rc,ChartId)
	if err != nil {
		middleware.GetLogEntry(r).Error(xerrs.Details(err, respond.ErrMaxStack))
		respond.Nay(w, r, code, err)
		return
	}
	respond.Yay(w, r, code, map[string]string{
		"message": "Data berhasil disimpan !",
	})
	return

}
