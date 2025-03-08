package request

type PurchaseDetailStore struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required"`
}

type PurchaseStore struct {
	PaymentMethod   string                `json:"payment_method" validate:"required"`
	PurchaseDetails []PurchaseDetailStore `json:"purchase_details" validate:"required"`
}

type QueryPurchase struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty"`
	UserId int    `validate:"omitempty"`
}
