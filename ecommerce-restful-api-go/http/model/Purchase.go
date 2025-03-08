package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Purchase struct {
	BaseModel
	UserID          uint64            `json:"user_id" gorm:"not null"`
	Code            string            `json:"code"`
	Total           float64           `json:"total"`
	PaymentMethod   string            `json:"payment_method"`
	Status          string            `json:"status"`
	User            *User             `gorm:"foreignKey:user_id;references:id"`
	DeletedAt       gorm.DeletedAt    `json:"-"`
	PurchaseDetails *[]PurchaseDetail `json:"purchase_details" gorm:"foreignKey:PurchaseId"`
}

func (purchase Purchase) GenerateCode(prefix string, db *gorm.DB) string {
	code := prefix
	codeLength := len(code)

	res := struct {
		Count uint64 `json:"count"`
	}{}

	resError := db.Model(purchase).Unscoped().Where("SUBSTRING(code,1,?) = ?", codeLength, code).
		Select("CONVERT(COALESCE(max(SUBSTRING(code,-7,7)),0),UNSIGNED) as count").
		First(&res)
	if resError.Error != nil {
		panic(resError.Error)
	}

	count := res.Count + 1
	code = fmt.Sprintf(code+"%07d", count)
	return code
}
