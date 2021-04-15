package model

type File struct {
	Model
	Name  string `gorm:"not null;size:32" json:"name"`
	Size  int    `gorm:"not null" json:"size"`
	Owner string `gorm:"not null;size:32" json:"owner"`
}
