package response

import "gorm.io/gorm"

type PurchaseResponse struct {
	Id              uint64                   `json:"id"`
	UserID          uint64                   `json:"user_id" gorm:"not null"`
	Code            string                   `json:"code"`
	Total           float64                  `json:"total"`
	PaymentMethod   string                   `json:"payment_method"`
	Status          string                   `json:"status"`
	User            *UserResponse            `gorm:"foreignKey:user_id;references:id"`
	DeletedAt       gorm.DeletedAt           `json:"-"`
	PurchaseDetails []PurchaseDetailResponse `json:"purchase_details" gorm:"foreignKey:PurchaseId;references:ID"`
}

func (PurchaseResponse) TableName() string {
	return "purchases"
}
