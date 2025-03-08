package model

import "gorm.io/gorm"

type User struct {
	BaseModel
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:user;not null" json:"role"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
