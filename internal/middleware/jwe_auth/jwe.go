package jwe_auth

import (
	"crypto/rsa"
	"time"

	"github.com/kopjenmbeng/goconf"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	ISSUER          = "onelabs.co"
	CfgJweOtpExpire = "jwe.otp_expire"

	DefaultExpiration    time.Duration = time.Duration(1) * time.Hour
	DefaultOtpExpiration time.Duration = time.Duration(5) * time.Minute
)

var (
	ErrInvalidScope = errors.New("invalid scope in claims")
	ErrNoTokenFound = errors.New("auth: no credentials attached in request")
)

type (
	JWE struct {
		i bool
		k *rsa.PrivateKey
		b jwt.NestedBuilder
		e int // expiration
	}
	PrivateClaims struct {
		DeviceId string   `json:"device_id"`
		Scopes   []string `json:"scopes,omitEmpty"`
		ClientId string   `json:"client_id,omitEmpty"`
	}
)

func (c PrivateClaims) Validate(e PrivateClaims) error {
	if len(e.Scopes) == 0 {
		return nil
	}
	if len(e.Scopes) == 0 {
		return ErrInvalidScope
	}
	sm := 0
	for _, s := range e.Scopes {
		for _, v := range c.Scopes {
			if v == s {
				sm++
			}
		}
	}
	if sm < len(e.Scopes) {
		return ErrInvalidScope
	}
	return nil
}

func NewJWE(key *rsa.PrivateKey, exp int) *JWE {
	return &JWE{k: key, e: exp}
}

func (j *JWE) builder() (builder jwt.NestedBuilder, err error) {
	if j.i {
		return j.b, nil
	}

	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: j.k}

	// create a Square.jose RSA signer, used to sign the JWT
	signerOpts := (&jose.SignerOptions{}).WithContentType("JWT")
	rsaSigner, err := jose.NewSigner(signingKey, signerOpts)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	encOpts := (&jose.EncrypterOptions{}).WithContentType("JWT")
	enc, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: &j.k.PublicKey}, encOpts)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	builder = jwt.SignedAndEncrypted(rsaSigner, enc)

	j.b = builder
	j.i = true

	return
}

func (j *JWE) Encode(pub jwt.Claims) (token string, err error) {
	var b jwt.NestedBuilder
	b, err = j.builder()
	if err != nil {
		return
	}

	if token, err = b.Claims(pub).CompactSerialize(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (j *JWE) Decode(token string) (pub jwt.Claims, err error) {
	var parsed *jwt.NestedJSONWebToken
	parsed, err = jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	decrypted, err := parsed.Decrypt(j.k)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if err = decrypted.Claims(&j.k.PublicKey, &pub); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (j *JWE) Valid() bool {
	return j.Valid()
}

func (j *JWE) GetExpire() time.Duration {
	if j.e < 1 {
		return DefaultExpiration
	}
	return time.Duration(j.e) * time.Hour
}

func (j *JWE) GetOtpExpire() time.Duration {
	otpExpire := goconf.GetInt(CfgJweOtpExpire)
	if otpExpire < 1 {
		return DefaultOtpExpiration
	}
	return time.Duration(otpExpire) * time.Minute
}
