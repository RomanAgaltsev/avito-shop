// Package auth provides methods for authorisation.
package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

type UserName string

const (
	// JWTSignAlgorithm contains JWT signing algorithm.
	JWTSignAlgorithm = "HS256"

	// UserNameClaimName contains key name of user name in a context.
	UserNameClaimName UserName = "username"
)

var ErrInvalidUser = fmt.Errorf("absent or invalid user in request")

// NewAuth returns new JWTAuth.
func NewAuth(secretKey string) *jwtauth.JWTAuth {
	return jwtauth.New(JWTSignAlgorithm, []byte(secretKey), nil)
}

// NewJWTToken creates new JWT token.
func NewJWTToken(ja *jwtauth.JWTAuth, username string) (token jwt.Token, tokenString string, err error) {
	return ja.Encode(map[string]interface{}{string(UserNameClaimName): username})
}

// HashPassword generates and returns hash of a given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares given password and hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UserFromRequest extracts user (name) from the given HTTP request.
func UserFromRequest(r *http.Request, secretKey string) (model.User, error) {
	// Create new JWT auth
	ja := NewAuth(secretKey)

	// Get JWT token from the cookie
	tokenString := jwtauth.TokenFromHeader(r)
	if tokenString == "" {
		return model.User{}, ErrInvalidUser
	}

	// Decode token string
	token, err := ja.Decode(tokenString)
	if err != nil {
		return model.User{}, err
	}

	// Get claims
	claims := token.PrivateClaims()

	// Get user name in interface type
	userNameInterface, ok := claims[string(UserNameClaimName)]
	if !ok {
		return model.User{}, ErrInvalidUser
	}

	// Convert user name to string
	userName, ok := userNameInterface.(string)
	if !ok {
		return model.User{}, ErrInvalidUser
	}

	// Return user structure
	return model.User{
		UserName: userName,
	}, nil
}
