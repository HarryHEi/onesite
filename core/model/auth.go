package model

// User 用户账户
type User struct {
	Model
	Username string `gorm:"unique;not null;size:32" json:"username,omitempty"`
	Password string `gorm:"not null;size:128" json:"password,omitempty"`
	Name     string `gorm:"not null;size:64" json:"name,omitempty"`
	IsAdmin  bool   `gorm:"not null" json:"is_admin,omitempty"`
}
