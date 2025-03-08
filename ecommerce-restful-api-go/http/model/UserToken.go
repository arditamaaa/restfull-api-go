package model

import (
	"time"
)

type UserToken struct {
	BaseModel
	Token   string     `gorm:"not null"`
	UserID  uint64     `json:"user_id" gorm:"not null"`
	Expires *time.Time `json:"expires"`
	User    *User      `gorm:"foreignKey:user_id;references:id"`
}
