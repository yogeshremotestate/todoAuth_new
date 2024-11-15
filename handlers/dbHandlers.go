package handlers

import (
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UserExist(c *gin.Context, email string) (user models.User, err error) {
	query := "SELECT id,email,password FROM users WHERE email = $1"

	err = initializers.DB.GetContext(c, &user, query, email)
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("userID", email),
		)
	}
	if err == sql.ErrNoRows {
		return user, errors.New("user not found")

	}
	return user, err

}

func CreateUser(c *gin.Context, email string, hash string) error {
	var user models.User
	query := `INSERT INTO users ( email, password,created_at,updated_at)
			  VALUES (TRIM($1), TRIM($2),$3,$4) RETURNING id`

	err := initializers.DB.Get(&user.ID, query, email, hash, time.Now(), time.Now())
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("userID", email),
			zap.String("hash", hash),
		)
	}
	return err
}

func CreateNote(c *gin.Context, title string, body string, userId uint) (err error) {
	var note models.Note
	query := `INSERT INTO notes ( title, body,user_id,created_at,updated_at)
			  VALUES ($1, $2,$3,$4,$5) RETURNING id`

	err = initializers.DB.Get(&note.ID, query, title, body, userId, time.Now(), time.Now())
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("title", title),
			zap.String("body", body),
		)
	}
	return err
}

func GetAll(c *gin.Context, userId uint) ([]models.Note, error) {
	var notes []models.Note
	query := `SELECT id, created_at, updated_at, title, body, user_id
        FROM notes 
        WHERE user_id = $1 and deleted_at IS NULL`
	err := initializers.DB.Select(&notes, query, userId)
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("query", query),
			zap.Uint("userId", userId),
		)
	}

	return notes, err
}

func GetOne(c *gin.Context, id string, userId uint) (models.Note, error) {
	var note models.Note
	query := "SELECT id,title,body, created_at,updated_at,user_id FROM notes WHERE id = $1 and user_id=$2"

	err := initializers.DB.GetContext(c, &note, query, id, userId)
	fmt.Println(err)
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("id", id),
			zap.Uint("userId", userId),
		)
	}

	return note, err
}

func UpdateOne(c *gin.Context, title string, body string, id string, userId int) (sql.Result, error) {
	updateQuery := `
        UPDATE notes 
        SET title = $1,body = $2
        WHERE id = $3 AND user_id = $4`
	result, err := initializers.DB.Exec(updateQuery, title, body, id, userId)
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("title", title),
			zap.String("body", body),
			zap.String("id", id),
			zap.Int("userId", userId),
		)
	}

	return result, err
}

func DeleteOne(c *gin.Context, id string, userId int) (sql.Result, error) {
	// log := middleware.GetLogger(c.Request.Context())
	query := `
    UPDATE notes 
    SET deleted_at = $1 
    WHERE id = $2 AND user_id = $3`
	result, err := initializers.DB.Exec(query, time.Now(), id, userId)
	if err != nil {
		zap.L().Warn("Executing SQL query",
			zap.String("id", id),
			zap.Int("userId", userId),
		)
	}

	return result, err
}
func ExcelRead(c *gin.Context, notes []models.Note) error {

	tx, err := initializers.DB.Beginx()
	if err != nil {
		zap.L().Info(err.Error())
		return err
	}

	for _, note := range notes {
		query := `INSERT INTO notes (title, body,user_id,created_at,updated_at) VALUES (:title, :body,:user_id,:created_at,:updated_at)`
		_, err = tx.NamedExec(query, note)
		if err != nil {
			tx.Rollback()
			zap.L().Warn(err.Error(),
				zap.String("query", query),
				zap.String("title", note.Title),
				zap.String("body", note.Body),
				zap.Uint("userId", note.UserID),
			)
			return err
		}
	}
	tx.Commit()
	return err
}
