package models

import "time"

type CommonFields struct {
	ID        uint       `gorm:"primaryKey, index" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}
