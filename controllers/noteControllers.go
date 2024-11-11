package controllers

import (
	"LearnGo-todoAuth/handlers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/middleware"
	"LearnGo-todoAuth/models"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoteCreate(c *gin.Context) {
	log := middleware.GetLogger(c.Request.Context())
	log.Info("NoteCreate is running")
	var body models.NoteBody
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}

	userDetail, _ := c.Get(initializers.UserString)
	err = handlers.CreateNote(c, body.Title, body.Body, uint(userDetail.(models.User).ID))
	if err == sql.ErrNoRows {
		c.JSON(400, gin.H{
			"error": errors.New("note not found"),
		})
		return
	} else if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Note addess successfully",
	})
}

func GetAllNote(c *gin.Context) {
	log := middleware.GetLogger(c.Request.Context())
	log.Info("GetAllNote is running")
	userDetail, _ := c.Get(initializers.UserString)

	notes, err := handlers.GetAll(c, uint(userDetail.(models.User).ID))
	if err == sql.ErrNoRows {
		c.JSON(400, gin.H{
			"error": errors.New("note not found"),
		})
		return
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"Posts": notes,
	})
}

func GetNote(c *gin.Context) {
	log := middleware.GetLogger(c.Request.Context())
	log.Info("GetOne is running")
	id := c.Param("id")
	userDetail, _ := c.Get("user")

	note, err := handlers.GetOne(c, id, uint(userDetail.(models.User).ID))
	if err == sql.ErrNoRows {
		c.JSON(400, gin.H{"error": errors.New("note not found")})
		return
	}
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{
		"Post": note,
	})
}

func UpdateNote(c *gin.Context) {
	log := middleware.GetLogger(c.Request.Context())
	log.Info("UpdateNote is running")
	id := c.Param("id")

	var body models.NoteBody
	Err := c.Bind(&body)
	if Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please send valid body"})
		log.Warn("request body is invalid")
		return
	}
	userDetail, _ := c.Get(initializers.UserString)

	_, err := handlers.GetOne(c, id, uint(userDetail.(models.User).ID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query note: %v", err)})
			return
		}
	}

	result, err := handlers.UpdateOne(c, body.Title, body.Body, id, userDetail.(models.User).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note Updated successfully"})
	log.Info("success")
}

func DeleteNote(c *gin.Context) {
	log := middleware.GetLogger(c.Request.Context())
	log.Info("DeleteNote is running")
	id := c.Param("id")
	userDetail, _ := c.Get(initializers.UserString)

	result, err := handlers.DeleteOne(c, id, userDetail.(models.User).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// rows affected check
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
