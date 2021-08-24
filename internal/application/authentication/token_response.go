package authentication

type TokenResponse struct {
	Token  string `json:"token"`
	Expiry int64    `json:"exp"`
}
