package controllers

import (
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"os"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(c *gin.Context) {
	// get body
	var body struct {
		Email    string
		Password string
	}
	err := c.Bind(&body) // pas the body from request

	//can apply mutiple checks on password for length and characters but keeping it simple for now
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}

	// fmt.Println(&body)

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	// fmt.Println(hash)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, result.Error)
		return
	}

	//return
	c.JSON(200, gin.H{
		"user": user, // return does not sends the new data
	})
}

func LoginUser(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	c.Bind(&body)

	//user exist check
	var user models.User
	initializers.DB.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": " email does not exist",
		})
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
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
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Generate token ",
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// Return JWT token in Response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

	// c.Status(200) // token return in cookie to this should be enough

}
