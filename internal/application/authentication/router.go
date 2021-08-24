package authentication

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/db_context"
)

const (
	CtxAuthseCaseKey = "auth_usecase"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(InjectUseCaseContext)
		r.Get("/get_token", GetTokenHandler)
		r.Post("/register", RegisterHandler)
	})
	return r
}

func InjectUseCaseContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbr := db_context.GetDbRead(r)
		dbw:=db_context.GetDbWrite(r)
		repo := NewAuthenticationRepository(dbr,dbw)
		usecase := NewAuthenticationUseCase(repo,r)
		ctx := context.WithValue(r.Context(), CtxAuthseCaseKey, usecase)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseCaseFromContext(c context.Context) IAuthenticationUseCase {
	return c.Value(CtxAuthseCaseKey).(IAuthenticationUseCase)
}
