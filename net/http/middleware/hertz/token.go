package hertz

import (
	"context"
	"crypto"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenHeader = "Authorization"
	TokenPrefix = "Bearer "
)

//type JWTOption struct {
//	PrivateKey crypto.PrivateKey
//	PublicKey  crypto.PublicKey
//}

type Jwt struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func (j *Jwt) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(j.PrivateKey)
}

func (j *Jwt) ParseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
}

func Token(j *Jwt) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tokenHeader := c.Request.Header.Get(TokenHeader)
		token, found := strings.CutPrefix(tokenHeader, TokenPrefix)
		if !found {
			return
		}

		fmt.Println(token)
	}
}
