package controllers

import (
	"LearnGo-todoAuth/handlers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"database/sql"
	"errors"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(c *gin.Context) {
	// log := middleware.GetLogger(c.Request.Context())
	zap.L().Info("SignUpUser is running")

	var body models.UserBody
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	Err := handlers.CreateUser(c, body.Email, string(hash))
	if Err == sql.ErrNoRows {
		c.JSON(400, gin.H{
			"error": errors.New("note not found"),
		})
		return
	} else if Err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"user": "User created successfully",
	})
}

func LoginUser(c *gin.Context) {
	zap.L().Info("LoginUser is running")
	var body models.UserBody
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}

	user, err := handlers.UserExist(c, body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	} else if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": " user not found",
		})
		return
	}

	Err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}
	// create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(initializers.ENV.SECRET))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Generate token ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
