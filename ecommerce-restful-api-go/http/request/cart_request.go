package request

type CartStore struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required"`
}

type CartPaymentStore struct {
	PaymentMethod string      `json:"payment_method" validate:"required"`
	CartDetails   []CartStore `json:"cart_details" validate:"required"`
}
