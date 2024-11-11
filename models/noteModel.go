package models

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        uint         `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"` // Optional for soft delete, if needed
	Title     string       `db:"title"`
	Body      string       `db:"body"`
	UserID    uint         `db:"user_id"`
}
