package dto

type TokenResponse struct {
	Token          string `json:"token"`
	ExpirationTime int64  `json:"expiration_time"`
}
