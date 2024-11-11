package middleware

import (
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Validate(c *gin.Context) {
	log := GetLogger(c.Request.Context())
	log.Info("Validating User i")
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Warn("Empty token")
		return
	}
	text := (strings.Split(bearerToken, " "))
	token, err := jwt.Parse(text[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(initializers.ENV.SECRET), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect token"})
		log.Warn("incorrect token")
		return

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		query := "SELECT id,email,password FROM users WHERE id = $1"
		err = initializers.DB.Get(&user, query, claims["userId"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Info("No user found")
		}
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
