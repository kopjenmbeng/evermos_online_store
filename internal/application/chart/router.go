package chart

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/db_context"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
)

const (
	CtxChartCaseKey = "chart_usecase"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(jwe_auth.GuardAnonymous(jwe_auth.TokenFromHeader))
		r.Use(InjectUseCaseContext)
		r.Post("/add", AddToChartHandler)
		r.Put("/update", UpdateChartHandler)
	})
	return r
}

func InjectUseCaseContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbr := db_context.GetDbRead(r)
		dbw := db_context.GetDbWrite(r)
		repo := NewChartRepository(dbr, dbw)
		usecase := NewChartUseCase(repo, r)
		ctx := context.WithValue(r.Context(), CtxChartCaseKey, usecase)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseCaseFromContext(c context.Context) IChartUseCase {
	return c.Value(CtxChartCaseKey).(IChartUseCase)
}