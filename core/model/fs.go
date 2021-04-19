package model

type File struct {
	Model
	Name     string `gorm:"not null;size:32" json:"name"`
	Fid      string `gorm:"not null;size:64" json:"fid"`
	Size     int    `gorm:"not null" json:"size"`
	Owner    string `gorm:"not null;size:32" json:"owner"`
	Exported bool   `gorm:"not null" json:"exported"`
}
