//go:generate mockgen -source=./jwt.go -destination=./mock/jwt.go -package=jwtmock
package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	key []byte
	exp time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
}

type Manager interface {
	// Issue a new token for a given id with exp time from env
	Issue(id string) (string, error)
	// Decode provided token to the uuid
	Decode(token string) (string, error)
}

// NewJWT - creates new instance of JWT.
func NewJWT(secret, validDays string) (*JWT, error) {
	atoi, err := strconv.Atoi(validDays)
	if err != nil {
		atoi = 14
	}

	return &JWT{
		key: []byte(secret),
		exp: time.Hour * 24 * time.Duration(atoi),
	}, err
}

// Issue - implementation of Manager.Issue.
func (j *JWT) Issue(id string) (string, error) {
	now := time.Now()
	exp := now.Add(j.exp)

	data := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString(j.key)
}

// Decode - implementation of Manager.Decode.
func (j *JWT) Decode(decode string) (string, error) {
	payload := &Claims{}
	_, err := jwt.ParseWithClaims(decode, payload, j.parseKey)
	if err != nil {
		return "", fmt.Errorf("token parse: %w", err)
	}

	if payload.Valid() != nil {
		return "", jwt.ErrTokenInvalidClaims
	}

	return payload.ID, nil
}

// parseKey - validates token signing method and alg, then return JWT.key.
func (j *JWT) parseKey(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return j.key, nil
}
