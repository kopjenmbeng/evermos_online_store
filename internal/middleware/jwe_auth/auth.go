package jwe_auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	// "github.com/koala-proptech/authentication/internal/middleware/jwe_auth"
	// "github.com/koala-proptech/authentication/internal/middleware/jwe_context"
	// "github.com/koala-proptech/authentication/internal/utility/respond"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	claimsKey       = "claims"
	JWTRetrievalKey = "access_token"
)

type (
	// TokenRetrieval function of retrieving auth token from request
	// it should extract token from request if any and return accordingly
	TokenRetrieval func(r *http.Request) string
	Claims         struct {
		Private PrivateClaims
		Public  jwt.Claims
	}
)

var ErrInvalidClient = errors.New("client not registered")

func introspect(r *http.Request, tr ...TokenRetrieval) (pub jwt.Claims, err error) {
	var ts string
	for _, fn := range tr {
		if ts = fn(r); ts != "" {
			break
		}
	}
	if ts == "" {
		err = ErrNoTokenFound
		return
	}

	return GetJWE(r.Context()).Decode(ts)
}

// GuardAnonymous middleware of allowing anonymous access to protected endpoints
// this will required request to contains proper JWE token even without proper subject
func GuardAnonymous(tr ...TokenRetrieval) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pub, err := introspect(r, tr...)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), claimsKey, &Claims{Public: pub})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}



func GetClaims(r *http.Request) (c *Claims) {
	if value := r.Context().Value(claimsKey); value != nil {
		c = value.(*Claims)
	}
	return
}

func GetMerchant(r *http.Request) (m string) {
	claims := GetClaims(r)
	if a := claims.Public.Audience; len(a) > 0 {
		m = a[0]
	}
	return
}

// TokenFromCookie retrieve authorization token from request cookie
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(JWTRetrievalKey)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// TokenFromHeader retrieve authorization token from request header
func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	if token := r.Header.Get("X-Auth-Token"); len(token) > 1 {
		return token
	}

	return ""
}

// TokenFromQuery retrieve authorization token from request query
func TokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get(JWTRetrievalKey)
}

// UserAgentFromHeader retrieve user agent from request header
func UserAgentFromHeader(r *http.Request) string {
	if agent := r.Header.Get("User-Agent"); len(agent) > 1 {
		return agent
	}
	return ""
}