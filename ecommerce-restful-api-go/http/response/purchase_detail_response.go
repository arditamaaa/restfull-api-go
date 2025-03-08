package response

type PurchaseDetailResponse struct {
	UserId     uint64          `json:"user_id"`
	ProductId  uint64          `json:"product_id"`
	Qty        int             `json:"qty"`
	TotalPrice float64         `json:"total_price"`
	User       UserResponse    `gorm:"foreignKey:UserId;references:User"`
	Product    ProductResponse `gorm:"foreignKey:ProductId;references:Product"`
}
