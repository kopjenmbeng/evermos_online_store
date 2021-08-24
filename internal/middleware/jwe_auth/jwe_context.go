package jwe_auth

import (
	"context"
	"net/http"

	// "github.com/koala-proptech/authentication/internal/middleware/jwe_auth"
)

const (
	JWEContextKey = "jwe"
)

func InitJWE(jw *JWE) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			if jw != nil {
				ctx = context.WithValue(r.Context(), JWEContextKey, jw)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func GetJWE(c context.Context) *JWE {
	return c.Value(JWEContextKey).(*JWE)
}
