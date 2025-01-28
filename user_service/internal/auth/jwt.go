package auth

import (
	"context"
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT struct {
	privateKey *rsa.PrivateKey

	accessToken  time.Duration
	refreshToken time.Duration
	issuer       string
	audience     []string
}

type AuthClaims struct {
	Role []string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTGenerator(certPath, issuer string, accessToken, refreshToken time.Duration, audience []string) (*JWT, error) {
	keyData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	return &JWT{privateKey, accessToken, refreshToken, issuer, audience}, nil
}

// TODO: а надо ли ? может возвращать структуру с access и refresh токеном
func (j *JWT) GenerateAllTokens(ctx context.Context, uid string, roles []string) ([]string, error) {

	return []string{}, nil
}

func (j *JWT) GenerateAccessToken(ctx context.Context, uid string, roles []string) (string, error) {
	claims := AuthClaims{
		Role: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(ctx, claims)

	return signedToken, err
}

func (j *JWT) GenerateRefreshToken(ctx context.Context, uid string, roles []string) (string, error) {
	claims := AuthClaims{
		Role: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(ctx, claims)

	return signedToken, err
}

func (j *JWT) generateJWT(ctx context.Context, claims AuthClaims) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JWT) generateJTI() string {
	return uuid.New().String()
}
