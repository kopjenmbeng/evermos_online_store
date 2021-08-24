package chart

import (
	"net/http"

	"github.com/RoseRocket/xerrs"
	"github.com/go-chi/render"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware"
	"github.com/kopjenmbeng/evermos_online_store/internal/utility/respond"
)

func AddToChartHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// err error
		rc  = r.Context()
		req AddChartRequest
		err error
	)
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	err = render.Bind(r, &req)
	if err != nil {
		middleware.GetLogEntry(r).Error(xerrs.Details(err, respond.ErrMaxStack))
		respond.Nay(w, r, http.StatusBadRequest, err)
		return
	}
	// app_code:=r.Header.Get("X-Client-id")

	useCase := UseCaseFromContext(rc)
	code, err := useCase.AddToChart(rc, req)
	if err != nil {
		middleware.GetLogEntry(r).Error(xerrs.Details(err, respond.ErrMaxStack))
		respond.Nay(w, r, code, err)
		return
	}
	respond.Yay(w, r, http.StatusCreated, map[string]string{
		"message": "Data berhasil disimpan !",
	})
	return

}
