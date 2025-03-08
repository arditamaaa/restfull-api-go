package response

type CartResponse struct {
	UserId     uint64          `json:"user_id" gorm:"not null"`
	ProductId  uint64          `json:"product_id" gorm:"not null"`
	Qty        int             `json:"qty"`
	TotalPrice float64         `json:"total_price"`
	User       UserResponse    `gorm:"-"`
	Product    ProductResponse `gorm:"-"`
}
