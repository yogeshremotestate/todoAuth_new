package controllers

import (
	"LearnGo-todoAuth/handlers"
	"LearnGo-todoAuth/initializers"
	Log "LearnGo-todoAuth/log"
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

// @Summary      Create Note
// @Description  Create a new note for the logged-in user
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note body models.NoteBody true "Create Note"
// @Security     BearerAuth
// @Success      200  "Success"
// @Router       /notes [post]
func NoteCreate(c *gin.Context) {
	// log := middleware.GetLogger(c.Request.Context())
	log := Log.GetLogger(c)
	log.Info("NoteCreate is running")
	var body models.NoteBody
	err := c.Bind(&body)
	if err != nil {
		log.Error(err.Error(), zap.String("err", "please send a valid body"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please send valid body",
		})
		return
	}
	log.Info("request Body", zap.String("title", body.Title), zap.String("body", body.Body))

	userDetail, _ := c.Get(initializers.UserString)
	err = handlers.CreateNote(c, body.Title, body.Body, uint(userDetail.(models.User).ID))
	if err == sql.ErrNoRows {
		log.Info("note not created")
		c.JSON(400, gin.H{
			"error": errors.New("note not found"),
		})
		return
	} else if err != nil {
		log.Info(err.Error())
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Note addess successfully",
	})
}

// @Summary      Get All Notes
// @Description  Retrieve all notes for the logged-in user
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Note "List of Notes"
// @Router       /notes [get]
func GetAllNote(c *gin.Context) {
	log := Log.GetLogger(c)
	zap.L().Info("GetAllNote is running")
	userDetail, _ := c.Get(initializers.UserString)

	notes, err := handlers.GetAll(c, uint(userDetail.(models.User).ID))
	if err == sql.ErrNoRows {
		log.Info("note not found")
		c.JSON(400, gin.H{
			"error": errors.New("note not found"),
		})
		return
	}
	if err != nil {
		log.Info(err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Posts": notes,
	})
}

// @Summary      Get Note
// @Description  Retrieve a specific note by its ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id path int true "Note ID"
// @Security     BearerAuth
// @Success      200 {object} models.Note "Retrieved Note"
// @Router       /notes/{id} [get]
func GetNote(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetOne is running")
	id := c.Param("id")
	// userDetail, _ := c.Get("user")

	note, err := handlers.GetOne(c, id)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		c.JSON(400, gin.H{"error": "note not found"})
		return
	}
	if err != nil {
		log.Info(err.Error())
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{
		"Post": note,
	})
}

// @Summary      Update Note
// @Description  Update an existing note by its ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id path int true "Note ID"
// @Param        note body models.NoteBody true "Update Note"
// @Security     BearerAuth
// @Success      200  "Success"
// @Router       /notes/{id} [put]
func UpdateNote(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("UpdateNote is running")
	id := c.Param("id")

	var body models.NoteBody
	Err := c.Bind(&body)
	if Err != nil {
		log.Info(Err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please send valid body"})
		return
	}
	log.Info("request Body", zap.String("title", body.Title), zap.String("body", body.Body))

	_, err := handlers.GetOne(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("Note not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		} else {
			log.Info(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query note: %v", err)})
			return
		}
	}

	result, err := handlers.UpdateOne(c, body.Title, body.Body, id)
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		log.Info("no row updated")
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note Updated successfully"})
	zap.L().Info("success")
}

// @Summary      Delete Note
// @Description  Delete a specific note by its ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id path int true "Note ID"
// @Security     BearerAuth
// @Success      200  "Success"
// @Router       /notes/{id} [delete]
func DeleteNote(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("DeleteNote is running")
	id := c.Param("id")
	// userDetail, _ := c.Get(initializers.UserString)

	result, err := handlers.DeleteOne(c, id)
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute update query: %v", err)})
		return
	}

	// rows affected check
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve rows affected: %v", err)})
		return
	}
	if rowsAffected == 0 {
		log.Info("No rows were updated")
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows were updated; check if the note exists with the given id and user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// @Summary      Upload Excel File
// @Description  Allows a logged-in user to upload an Excel file
// @Tags         file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Excel file to upload"
// @Success      200 "Success"
// @Security     BearerAuth
// @Router       /notes/upload [post]
func UploadExcel(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("UploadExcel is running")
	file, err := c.FormFile("file")
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	filePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save error"})
		return
	}

	notes, err := services.ParseNotesFromExcel(c, filePath)
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing Excel"})
		return
	}

	err = handlers.ExcelRead(c, notes)
	if err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notes uploaded successfully"})

}
