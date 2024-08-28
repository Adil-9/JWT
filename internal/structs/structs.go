package structs

import "github.com/golang-jwt/jwt/v5"

type Data struct {
	IP string `json:"ip"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	Hash   []byte
	IP     string
}