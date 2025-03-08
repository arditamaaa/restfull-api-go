package request

type ProductStore struct {
	Name  string  `json:"name" validate:"required,max=50" example:"fake name"`
	Price float64 `json:"price" validate:"required"`
}

type QueryProduct struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}
