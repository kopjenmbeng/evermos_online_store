package jwe_auth

import (
	"net/http"
	"time"

	chim "github.com/go-chi/chi/middleware"
	// "github.com/koala-proptech/authentication/internal/middleware/jwe_auth"
	"gopkg.in/square/go-jose.v2/jwt"
)

func GenerateToken(r *http.Request, userId string, isOtp bool) (token string, expiry int64, err error) {
	t := time.Now()
	jw := GetJWE(r.Context())

	timeExpiry := jw.GetExpire()
	if isOtp {
		timeExpiry = jw.GetOtpExpire()
	}

	pub := jwt.Claims{
		Issuer:    ISSUER,
		// Audience:  jwt.Audience{merchant},
		ID:        chim.GetReqID(r.Context()),
		Subject:   userId,
		Expiry:    jwt.NewNumericDate(t.Add(timeExpiry)),
		IssuedAt:  jwt.NewNumericDate(t),
		NotBefore: jwt.NewNumericDate(t),
	}
	// pri := PrivateClaims{
	// 	DeviceId: deviceId,
	// 	Scopes:   scopes,
	// 	ClientId: clientId,
	// }

	token, err = jw.Encode(pub)
	return token, int64(*pub.Expiry), err
}

func RefreshToken(r *http.Request, claims *Claims, expiryAt int64) (string, error) {
	jw := GetJWE(r.Context())
	claims.Public.Expiry = jwt.NewNumericDate(time.Unix(expiryAt, 0))
	return jw.Encode(claims.Public)
}