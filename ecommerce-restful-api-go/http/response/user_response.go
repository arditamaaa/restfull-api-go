package response

type UserResponse struct {
	Id    uint64 `json:"Id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (UserResponse) TableName() string {
	return "users"
}
