package controllers

import (
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoteCreate(c *gin.Context) {
	// get body
	var body struct {
		Title string
		Body  string
	}
	err := c.Bind(&body) // pas the body from request

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}
	user, _ := c.Get("user")

	// create Note
	post := models.Note{Title: body.Title, Body: body.Body, UserID: user.(models.User).ID}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(http.StatusBadGateway)
		return
	}

	//return
	c.JSON(200, gin.H{
		"Post": post, // can not return result var as it does not show anything on response
	})
}

func GetAllNote(c *gin.Context) {

	// get all Notes
	var posts []models.Note

	check, _ := c.Get("user")

	initializers.DB.Where("user_id=?", check.(models.User).ID).Preload("User").Find(&posts)

	//return
	c.JSON(200, gin.H{
		"Posts": posts,
	})
}

func GetNote(c *gin.Context) {
	//get id
	id := c.Param("id")
	// fmt.Println(id)
	user, _ := c.Get("user")

	// get all Notes
	var post models.Note

	// initializers.DB.First(&post, id[0].Value)
	initializers.DB.Where("user_id=?", user.(models.User).ID).Preload("User").First(&post, "id=?", id)

	//return
	c.JSON(200, gin.H{
		"Post": post,
	})
}

func UpdateNote(c *gin.Context) {
	// get body and id
	id := c.Param("id")

	var body struct {
		Title string
		Body  string
	}
	err := c.Bind(&body) // pass the body

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}
	user, _ := c.Get("user")

	// find and update
	var post models.Note

	initializers.DB.Where("user_id=?", user.(models.User).ID).Preload("User").First(&post, "id=?", id)

	result := initializers.DB.Model(&post).Updates(models.Note{Title: body.Title, Body: body.Body})

	if result.Error != nil {
		c.Status(400)
		return
	}

	//return
	c.JSON(200, gin.H{
		"Post": post, // can not return result var as it does not show anything on response
	})
}

func DeleteNote(c *gin.Context) {
	//get id
	id := c.Params
	// fmt.Println(id)
	user, _ := c.Get("user")

	// we can first search the note and compare the user id before deleting

	// for now just doing delete, no check
	initializers.DB.Where("user_id=?", user.(models.User).ID).Delete(&models.Note{}, id[0].Value)

	//return
	c.Status(200)
}
