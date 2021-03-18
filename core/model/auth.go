package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:32"`
	Password string `gorm:"not null;size:128"`
	Name     string `gorm:"not null;size:64"`
	IsAdmin  bool   `gorm:"not null"`
}
