package auth

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobwas/glob"
	"time"
)

type JWT struct {
	Account string `json:"account"`
	Prefix  string `json:"prefix"`
	Type    string `json:"type"`
	jwt.StandardClaims
}

type LocalJWTProvider struct{}

func (l LocalJWTProvider) NewJWT(account, prefix, types string, key []byte, ttlInSeconds int64) (string, error) {

	claim := JWT{
		account,
		prefix,
		types,
		jwt.StandardClaims{
			IssuedAt:  time.Now().UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
			ExpiresAt: time.Now().UTC().Add(time.Duration(ttlInSeconds) * time.Second).Unix(),
			Issuer:    "test",
		},
	}

	tokenMaterial := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenMaterial.SignedString(key)
	fmt.Printf("%v %v", token, err)

	return token, err
}

func (l LocalJWTProvider) VerifyJWT(base64JWT string, key []byte) (JWT, error) {

	token, err := jwt.ParseWithClaims(base64JWT, &JWT{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return JWT{}, err
	}

	if claim, ok := token.Claims.(*JWT); ok && token.Valid {
		return *claim, nil
	}

	return JWT{}, errors.New("Token Cannot Be Parsed or is Invalid")
}

func (l LocalJWTProvider) Authorize(jwt JWT, resource string) (bool, error) {
	g := glob.MustCompile(jwt.Prefix)

	matched := g.Match(resource)

	return matched, nil
}
