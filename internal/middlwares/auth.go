package middlewares

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type keyUserID string

const (
	ContextUserID keyUserID = "ID"
	NameCookie    string    = "x-auth-key"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &jwt.StandardClaims{}
		authKey := c.GetHeader("Authorization")

		if authKey == "" {
			authKey, err := setToken(uuid.NewString(), jwtSecret, claims)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)

				return
			}

			c.SetCookie(NameCookie, claims.Id, 86000, "/", c.Request.URL.Hostname(), true, true)
			c.Set(string(ContextUserID), claims.Id)
			c.Header("Authorization", authKey)
			c.Next()

			return
		}

		if err := validToken(authKey, jwtSecret, claims); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		userCookie, err := c.Cookie(NameCookie)
		if err != nil {
			c.SetCookie(NameCookie, claims.Id, 86000, "/", c.Request.URL.Hostname(), true, true)

			return
		}
		fmt.Println(userCookie)

		c.Set(string(ContextUserID), claims.Id)
		c.Next()
	}
}

func setToken(uid, jwtSecret string, claims *jwt.StandardClaims) (string, error) {
	claims.Id = uid
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create token for user %s", uid)
	}

	return ss, nil
}

func validToken(userToken, jwtSecret string, claims *jwt.StandardClaims) error {
	token, err := jwt.ParseWithClaims(userToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("token not valid")
	}

	return nil
}
