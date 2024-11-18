package controllers

import (
	"LearnGo-todoAuth/handlers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"database/sql"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      SignUpUser
// @Description  Register User new account
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        credentials body models.UserBody true "Signup Credentials"
// @Success      201 {object} map[string]interface{} "User Created Successfully"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      409 {object} map[string]string "Conflict - Email Already Exists"
// @Router       /user/signup [post]
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

	user, err := handlers.UserExist(c, body.Email)
	if user.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": " user email already exist",
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
			"error": "user not created",
		})
		return
	} else if Err != nil {
		c.JSON(400, gin.H{
			"error": Err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"user": "User created successfully",
	})
}

// @Summary      LoginUser
// @Description  Authenticate a user and return a token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body models.UserBody true "Login Credentials"
// @Success      200 {object} map[string]interface{} "Token"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /user/login [post]
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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " user not found",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
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
