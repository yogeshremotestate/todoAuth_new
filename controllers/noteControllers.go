package controllers

import (
	"LearnGo-todoAuth/handlers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"LearnGo-todoAuth/services"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NoteCreate(c *gin.Context) {
	// log := middleware.GetLogger(c.Request.Context())
	zap.L().Info("NoteCreate is running")
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
	zap.L().Info("GetAllNote is running")
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
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Posts": notes,
	})
}

func GetNote(c *gin.Context) {
	zap.L().Info("GetOne is running")
	id := c.Param("id")
	// userDetail, _ := c.Get("user")

	note, err := handlers.GetOne(c, id)
	if err == sql.ErrNoRows {
		zap.L().Info(err.Error())
		c.JSON(400, gin.H{"error": "note not found"})
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
	zap.L().Info("UpdateNote is running")
	id := c.Param("id")

	var body models.NoteBody
	Err := c.Bind(&body)
	if Err != nil {
		zap.L().Info(Err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please send valid body"})
		return
	}

	_, err := handlers.GetOne(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query note: %v", err)})
			return
		}
	}

	result, err := handlers.UpdateOne(c, body.Title, body.Body, id)
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note Updated successfully"})
	zap.L().Info("success")
}

func DeleteNote(c *gin.Context) {
	zap.L().Info("DeleteNote is running")
	id := c.Param("id")
	// userDetail, _ := c.Get(initializers.UserString)

	result, err := handlers.DeleteOne(c, id)
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// rows affected check
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func UploadExcel(c *gin.Context) {
	zap.L().Info("UploadExcel is running")
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	filePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save error"})
		return
	}

	notes, err := services.ParseNotesFromExcel(c, filePath)
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing Excel"})
		return
	}

	err = handlers.ExcelRead(c, notes)
	if err != nil {
		zap.L().Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notes uploaded successfully"})

}
