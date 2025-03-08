package model

import (
	"time"
)

type BaseModel struct {
	ID        uint64     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Tabler interface {
	TableName() string
}
