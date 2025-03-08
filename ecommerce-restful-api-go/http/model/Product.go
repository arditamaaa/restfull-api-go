package model

import (
	"gorm.io/gorm"
)

type Product struct {
	BaseModel
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
