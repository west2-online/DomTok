/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type Claims struct {
	Type   int64 `json:"type"`
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// CreateAllToken 创建一对 token，第一个是 access token，第二个是 refresh token
func CreateAllToken(uid int64) (string, string, error) {
	accessToken, err := CreateToken(constants.TypeAccessToken, uid)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := CreateToken(constants.TypeRefreshToken, uid)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// CreateToken 根据 token 类型和用户 ID 创建 token
func CreateToken(tokenType int64, uid int64) (string, error) {
	if config.Server == nil {
		return "", errno.NewErrNo(errno.AuthInvalidCode, "server config not found")
	}

	expireTime := time.Now().Add(getTokenTTL(tokenType))
	claims := Claims{
		Type:   tokenType,
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    constants.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	key, err := parsePrivateKey(config.Server.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// CheckToken 检查 token 是否有效，并返回 token 类型和用户 ID
func CheckToken(token string) (int64, int64, error) {
	if config.Server == nil {
		return 0, 0, errno.NewErrNo(errno.AuthInvalidCode, "server config not found")
	}
	if token == "" {
		return -1, 0, errno.NewErrNo(errno.AuthMissingTokenCode, "token is empty")
	}

	unverifiedClaims, err := parseUnverifiedClaims(token)
	if err != nil {
		return -1, 0, fmt.Errorf("failed to parse unverified claims: %w", err)
	}

	secret, err := parsePublicKey(config.Server.PublicKey)
	if err != nil {
		return -1, 0, fmt.Errorf("failed to parse public key: %w", err)
	}

	verifiedClaims, err := verifyToken(token, secret)
	if err != nil {
		tokenType, err := handleTokenError(err, unverifiedClaims.Type)
		return tokenType, 0, err
	}

	return verifiedClaims.Type, verifiedClaims.UserID, nil
}

// parsePrivateKey 解析 Ed25519 私钥
func parsePrivateKey(key string) (interface{}, error) {
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return privateKey, nil
}

// parsePublicKey 解析 Ed25519 公钥
func parsePublicKey(key string) (interface{}, error) {
	publicKey, err := jwt.ParseEdPublicKeyFromPEM([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return publicKey, nil
}

// parseUnverifiedClaims 解析未验证的 token claims
func parseUnverifiedClaims(token string) (*Claims, error) {
	tokenStruct, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse unverified token: %w", err)
	}

	claims, ok := tokenStruct.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	return claims, nil
}

// verifyToken 验证 token 并返回 claims
func verifyToken(token string, key interface{}) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// handleTokenError 处理 token 验证错误
func handleTokenError(err error, tokenType int64) (int64, error) {
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			if tokenType == constants.TypeAccessToken {
				return -1, errno.AuthAccessExpired
			} else if tokenType == constants.TypeUserLoginToken {
				return -1, errno.AuthAccessExpired
			}
			return -1, errno.NewErrNo(errno.AuthRefreshExpiredCode, "refresh token expired")
		}
	}
	return -1, fmt.Errorf("token validation failed: %w", err)
}

// getTokenTTL 根据 token 类型返回过期时间
func getTokenTTL(tokenType int64) time.Duration {
	switch tokenType {
	case constants.TypeAccessToken:
		return constants.AccessTokenTTL
	case constants.TypeRefreshToken:
		return constants.RefreshTokenTTL
	case constants.TypeUserLoginToken:
		return constants.AccessTokenTTL
	default:
		return 0
	}
}
