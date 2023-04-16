package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("iamtheboneofmysword")

// MyClaims Customer jwt.StandardClaims
type MyClaims struct {
	Account string `json:"account"`
	jwt.StandardClaims
}

// GenToken Create a new token
func GenToken(account string) (string, error) {
	c := MyClaims{
		account,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    account,
		},
	}
	// Choose specific algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Choose specific Signature
	return token.SignedString(SecretKey)
}

// ParseToken Parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// Valid token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// code -2 jwt驗證未通過
// JWTAuthMiddleware Middleware of JWT
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		// Get token from Header.Authorization field.
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -2,
				"msg":  "Authorization is null in Header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -2,
				"msg":  "Format of Authorization is wrong",
			})
			c.Abort()
			return
		}
		// parts[0] is Bearer, parts is token.
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -2,
				"msg":  "Invalid Token.",
			})
			c.Abort()
			return
		}
		// Store Account info into Context
		c.Set("account", mc.Account)
		// After that, we can get Account info from c.Get("account")
		c.Next()
	}
}
