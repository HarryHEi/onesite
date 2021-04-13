package model

import (
	"time"

	"gorm.io/gorm"
)

// Model 数据库表通用结构
type Model struct {
	ID        uint            `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
