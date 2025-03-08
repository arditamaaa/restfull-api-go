package controller

import (
	"math"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService  service.UserService
	TokenService service.TokenService
}

func NewUserController(userService service.UserService, tokenService service.TokenService) *UserController {
	return &UserController{
		UserService:  userService,
		TokenService: tokenService,
	}
}

func (u *UserController) Index(c *fiber.Ctx) error {
	query := &request.QueryUser{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	users, totalResults, err := u.UserService.GetUsers(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.User]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all users successfully",
			Results:      users,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

func (u *UserController) Create(c *fiber.Ctx) error {
	req := new(request.StoreUser)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.CreateUser(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.Success[model.User]{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Created user successfully",
			Data:    *user,
		})
}

func (u *UserController) Show(c *fiber.Ctx) error {
	userID := c.Params("id")
	user, err := u.UserService.GetUserByID(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.User]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get user successfully",
			Data:    *user,
		})

}

func (u *UserController) Update(c *fiber.Ctx) error {
	req := new(request.UpdateUser)
	userID := c.Params("id")

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.UpdateUser(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.User]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update user successfully",
			Data:    *user,
		})
}

func (u *UserController) Delete(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := u.TokenService.DeleteToken(c, userID); err != nil {
		return err
	}

	if err := u.UserService.DeleteUser(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete user successfully",
		})
}
