package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/graytonio/flagops-data-storage/internal/services/user"
)

// Handles the parsing and validation of JWT authentication tokens
type JWTService struct {
	AccessExpires time.Duration
	RefreshExpires time.Duration
	SigningSecret string

	UserDataService *user.UserDataService
}

type UserClaims struct {
	ID          uint     `json:"id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type UserRefreshClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func (jh *JWTService) NewUserAccessToken(claims *UserClaims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jh.AccessExpires)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(jh.SigningSecret))
}

func (jh *JWTService) ParseUserAccessToken(accessToken string) (*UserClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jh.SigningSecret), nil
	})
	if err != nil {
	  return nil, err
	}

	claims, ok := parsedToken.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("mailformed token payload")
	}

	return claims, nil
}

func (jh *JWTService) NewUserRefreshToken(claims *UserRefreshClaims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jh.RefreshExpires)),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(jh.SigningSecret))
}

func (jh *JWTService) ParseUserRefreshToken(refreshToken string) (*UserRefreshClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &UserRefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jh.SigningSecret), nil
	})
	if err != nil {
	  return nil, err
	}

	claims, ok := parsedToken.Claims.(*UserRefreshClaims)
	if !ok {
		return nil, errors.New("mailformed token payload")
	}

	return claims, nil
}

// Get the user object from the access token. If access token needs to be refreshed it will be returned
func (jh *JWTService) ValidateUserTokens(accessToken string, refreshToken string) (*UserClaims, string, error) {
	parsedRefresh, refreshErr := jh.ParseUserRefreshToken(refreshToken)
	if refreshErr != nil { // Invalid refresh token can't renew access token
		return nil, "", refreshErr
	}

	if accessToken == "" {
		return jh.refreshAccessToken(parsedRefresh)
	}
	
	parsedAccess, accessErr := jh.ParseUserAccessToken(accessToken)
	if accessErr == nil { // Access token is valid fetch 
		return parsedAccess, "", nil
	}

	if !errors.Is(accessErr, jwt.ErrTokenExpired) { // Some other jwt parsing error
		return nil, "", accessErr 
	}

	// Refresh access token
	return jh.refreshAccessToken(parsedRefresh)
}

// Generates a new access token from the refresh token data
func (jh *JWTService) refreshAccessToken(refreshToken *UserRefreshClaims) (*UserClaims, string, error) {
	user, err := jh.UserDataService.GetUserByID(refreshToken.ID)
	if err != nil {
	  return nil, "", err
	}

	permissions := []string{}
	for _, p := range user.Permissions {
		permissions = append(permissions, p.ID)
	}

	claims := &UserClaims{
		ID: user.ID,
		Permissions: permissions,
	}

	newAccessToken, err := jh.NewUserAccessToken(claims)
	return claims, newAccessToken, err
}