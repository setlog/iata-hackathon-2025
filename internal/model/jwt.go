package model

import "time"

type JwtToken struct {
	Token              string
	TokenExpiry        time.Time
	RefreshToken       string
	RefreshTokenExpiry time.Time
}
