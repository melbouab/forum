package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func SignToken(userid, username string) (string, error) {
	jwtsecret := os.Getenv("JWT_SECRET")
	jwtexpirestime := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userid,
		"user": username,
	}

	if jwtexpirestime != "" {
		duration, err := time.ParseDuration(jwtexpirestime)
		if err != nil {
			return "", err
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signToken, err := token.SignedString([]byte(jwtsecret))
	if err != nil {
		return "", err
	}

	return signToken, nil
}
