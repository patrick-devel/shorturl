package middlewares

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type keyUserID string

const (
	ContextUserID keyUserID = "ID"
)

func AuthMiddleware(jwtSecret string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &jwt.StandardClaims{}

		var authToken string

		authHeader := c.GetHeader("Authorization")
		cookie, err := c.Cookie("user_id_signed")
		if err != nil {
			logrus.WithError(err).Warning("cookie not found")
		}

		switch {
		case authHeader != "":
			authToken = authHeader
		case cookie != "":
			authToken = cookie
		}

		logger.WithFields(logrus.Fields{"Authorization": authHeader, "cookie": cookie}).Info("check authorization")

		if authToken == "" {
			if c.Request.Method == http.MethodPost {
				authKey, err := setToken(uuid.NewString(), jwtSecret, claims)
				if err != nil {
					logger.Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)

					return
				}
				c.Header("Authorization", authKey)
				c.SetCookie("user_id_signed", authKey, 60*60*24, "/", "", false, true)
				c.Set(string(ContextUserID), claims.Id)
				c.Next()

				return
			}
		}

		if err := validToken(authToken, jwtSecret, claims); err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Set(string(ContextUserID), claims.Id)
		logger.WithField("claims", claims).Info("Info CLAIMS")
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
