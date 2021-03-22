package model

type User struct {
	Model
	Username string `gorm:"unique;not null;size:32" json:"username"`
	Password string `gorm:"not null;size:128" json:"password"`
	Name     string `gorm:"not null;size:64" json:"name"`
	IsAdmin  bool   `gorm:"not null" json:"is_admin"`
}
