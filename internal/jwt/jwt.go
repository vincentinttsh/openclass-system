package jwt

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var signingKey []byte

// MapClaims is a claims type that uses the map[string]interface{} for JSON decoding.
// This is the default claims type if you don't supply one
type MapClaims = jwt.MapClaims

func init() {
	signingKey = []byte(os.Getenv("JWT_SECRET"))
}

// Sign is a function that returns the signing key
func Sign(claims jwt.MapClaims) (t string, err error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	t, err = token.SignedString(signingKey)

	return
}

// Verify is a function that verifies the token
func Verify(tokenString string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	return
}
