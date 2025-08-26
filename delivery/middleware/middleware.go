package middleware

import (
	"github.com/teuflischone/my-go-api/utils/jwt"
)

// Middleware ...
type Middleware struct {
	jwtSvc jwt.JWTService
}

// NewMiddleware will create new Middleware object
func NewMiddleware(jwtSvc jwt.JWTService) *Middleware {
	return &Middleware{
		jwtSvc: jwtSvc,
	}
}
