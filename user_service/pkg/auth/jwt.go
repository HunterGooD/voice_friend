package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type JWT struct {
	privateKey *rsa.PrivateKey

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
	audience             []string
}

type AuthClaims struct {
	Role string `json:"role"`
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

func NewJWTGeneratorWithPrivateKey(privateKey *rsa.PrivateKey, issuer string, accessToken, refreshToken time.Duration, audience []string) (*JWT, error) {
	return &JWT{privateKey, accessToken, refreshToken, issuer, audience}, nil
}

func loadPublicKey(filePath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func ValidateJWT(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

// GenerateAllTokensAsync TODO: а надо ли ? может возвращать структуру с access и refresh токеном|
//
//	return array tokens first elem is access token and second if refresh token
func (j *JWT) GenerateAllTokensAsync(ctx context.Context, uid, role string) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	refreshCh, accessCh := make(chan string, 1), make(chan string, 1)
	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		access, err := j.GenerateAccessToken(ctx, uid, role)
		if err != nil {
			errCh <- errors.Wrap(err, "Error create access token")
			return
		}
		accessCh <- access
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		refresh, err := j.GenerateRefreshToken(ctx, uid, role)
		if err != nil {
			errCh <- errors.Wrap(err, "Error create refresh token")
			return
		}
		refreshCh <- refresh
	}()

	go func() {
		wg.Wait()
		close(errCh)
		close(accessCh)
		close(refreshCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return []string{<-accessCh, <-refreshCh}, nil
}

// GenerateAllTokens TODO: а надо ли ? может возвращать структуру с access и refresh токеном|
//
// return array tokens first elem is access token and second if refresh token
func (j *JWT) GenerateAllTokens(ctx context.Context, uid, role string) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	access, err := j.GenerateAccessToken(ctx, uid, role)
	if err != nil {
		return nil, errors.Wrap(err, "Error create access token")
	}

	refresh, err := j.GenerateRefreshToken(ctx, uid, role)
	if err != nil {
		return nil, errors.Wrap(err, "Error create refresh token")
	}

	return []string{access, refresh}, nil
}

func (j *JWT) GenerateAccessToken(ctx context.Context, uid, role string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	claims := AuthClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(&claims)

	return signedToken, err
}

func (j *JWT) GenerateRefreshToken(ctx context.Context, uid, role string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	claims := AuthClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(&claims)

	return signedToken, err
}

func (j *JWT) generateJWT(claims *AuthClaims) (string, error) {
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
