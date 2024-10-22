package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	_ "github.com/joho/godotenv/autoload"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaim struct {
	//UntisName string `json:"untisName"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string, persistantLogin bool) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	if persistantLogin {
		expirationTime = time.Now().Add(2160 * time.Hour)
	}
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtSecret)
	return
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, nil
}
