package response

type ProductResponse struct {
	Id    uint64  `json:"Id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (ProductResponse) TableName() string {
	return "products"
}
