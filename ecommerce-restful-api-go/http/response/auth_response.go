package response

import "simple-api-go/http/model"

type SuccessWithTokens struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
	Token   string     `json:"token"`
}
