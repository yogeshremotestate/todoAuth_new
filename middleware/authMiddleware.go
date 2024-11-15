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
	"go.uber.org/zap"
)

func AuthValidate(c *gin.Context) {
	// log := GetLogger(c.Request.Context())
	zap.L().Info("Validating User")
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		zap.L().Warn("Empty token")
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
		zap.L().Warn("incorrect token")
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
			zap.L().Info("No user found")
		}
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
func VerifyUserNote(c *gin.Context) {
	// log := GetLogger(c.Request.Context())
	zap.L().Info("verifying Note belongs to user")

	noteId := c.Param("id")
	userDetail, _ := c.Get("user")

	var note models.Note
	query := "SELECT id,title,body, created_at,updated_at,user_id FROM notes WHERE id = $1 "

	err := initializers.DB.GetContext(c, &note, query, noteId)
	if err != nil {
		zap.L().Info(err.Error(),
			zap.String("query", query),
			zap.String("id", noteId),
		)
		c.AbortWithStatus(http.StatusBadRequest)
		zap.L().Warn("query failed at execution")
		return
	}

	if note.UserID != uint(userDetail.(models.User).ID) {

		zap.L().Warn("Note does not belongs to user")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Note does not belongs to user"})
		return
	}

}
