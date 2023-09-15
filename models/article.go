package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Email     string         `json:"email"`
	Tittle    string         `json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
