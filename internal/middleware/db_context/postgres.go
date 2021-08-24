package db_context

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const (
	dbReadContextKey  = "dbr"
	dbWriteContextKey = "dbw"
)

func DBR(dbR *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			if dbR != nil {
				ctx = context.WithValue(ctx, dbReadContextKey, dbR)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func DBW(dbW *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			if dbW != nil {
				ctx = context.WithValue(ctx, dbWriteContextKey, dbW)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func GetDbRead(r *http.Request) *sqlx.DB {
	return r.Context().Value(dbReadContextKey).(*sqlx.DB)
}

func GetDbWrite(r *http.Request) *sqlx.DB {
	return r.Context().Value(dbWriteContextKey).(*sqlx.DB)
}
