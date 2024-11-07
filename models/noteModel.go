package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title  string
	Body   string
	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
