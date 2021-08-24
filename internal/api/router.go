package api

import (
	"github.com/go-chi/chi"
	"github.com/kopjenmbeng/evermos_online_store/internal/application/authentication"
	"github.com/kopjenmbeng/evermos_online_store/internal/application/chart"
	// "github.com/kopjenmbeng/evermos_online_store/internal/application/tes"
)

func routes(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/authentication", authentication.Routes())
		r.Mount("/chart", chart.Routes())
		// r.Mount("/socket", mysocket.Routes())
		// r.Mount("/auth", auth.Routes())
		// r.Mount("/prime", prime.Routes())
	})
}
