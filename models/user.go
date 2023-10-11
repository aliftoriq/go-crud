package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}
