package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mini-blog/internal/config"
	"mini-blog/pkg/apperror"
)

type TokenManager struct {
	secret             []byte
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

type JwtToken struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func NewTokenManager(auth config.Auth) *TokenManager {
	return &TokenManager{
		secret:             []byte(auth.JwtSecret),
		accessTokenExpire:  time.Duration(auth.AccessTokenExpireMins) * time.Minute,
		refreshTokenExpire: time.Duration(auth.RefreshTokenExpireDays) * time.Hour * 24,
	}
}

func (tm *TokenManager) Generate(userID int64) (JwtToken, error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tm.accessTokenExpire).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString(tm.secret)
	if err != nil {
		return JwtToken{}, apperror.ErrInternal
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tm.refreshTokenExpire).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString(tm.secret)
	if err != nil {
		return JwtToken{}, apperror.ErrInternal
	}
	return JwtToken{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (tm *TokenManager) Parse(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return tm.secret, nil
	})

	if err != nil || !token.Valid {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("%w: token is expire", apperror.ErrUnauthorized)
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, fmt.Errorf("%w: token is not active yet", apperror.ErrUnauthorized)
		default:
			return nil, fmt.Errorf("%w: invalid token", apperror.ErrForbidden)
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%w: invalid claims", apperror.ErrForbidden)
	}
	return claims, nil
}
