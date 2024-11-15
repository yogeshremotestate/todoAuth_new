package services

import (
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

func ParseNotesFromExcel(c *gin.Context, filePath string) ([]models.Note, error) {

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		zap.L().Info(err.Error())
		return nil, err
	}

	var notes []models.Note
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		zap.L().Info(err.Error())
		return nil, err
	}

	userDetail, _ := c.Get(initializers.UserString)

	for _, row := range rows[1:] { // Skip header
		if len(row) >= 2 {
			note := models.Note{
				Title:     row[0],
				Body:      row[1],
				UserID:    uint(userDetail.(models.User).ID),
				UpdatedAt: time.Now(),
				CreatedAt: time.Now(),
			}
			notes = append(notes, note)
		}
	}
	return notes, nil
}
