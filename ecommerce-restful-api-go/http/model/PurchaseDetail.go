package model

import "gorm.io/gorm"

type PurchaseDetail struct {
	BaseModel
	PurchaseId uint64         `json:"purchase_id" gorm:"not null"`
	ProductId  uint64         `json:"product_id" gorm:"not null"`
	Qty        int            `json:"qty"`
	DeletedAt  gorm.DeletedAt `json:"-"`
	Product    *Product       `json:"product" gorm:"foreignKey:product_id;references:id"`
}
