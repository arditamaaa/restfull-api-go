package request

type StoreUser struct {
	Name     string `json:"name" validate:"required,max=50" example:"fake name"`
	Email    string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
	Role     string `json:"role" validate:"required,oneof=user admin,max=50" example:"user"`
}

type UpdateUser struct {
	Name  string `json:"name,omitempty" validate:"omitempty,max=50" example:"fake name"`
	Email string `json:"email" validate:"omitempty,email,max=50" example:"fake@example.com"`
	Role  string `json:"role" validate:"required,oneof=user admin,max=50" example:"user"`
}

type QueryUser struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}
