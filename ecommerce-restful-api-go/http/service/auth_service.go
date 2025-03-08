package service

import (
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/validation"
	"simple-api-go/util"

	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(c *fiber.Ctx, req *request.Register) (*model.User, error)
	Login(c *fiber.Ctx, req *request.Login) (*model.User, error)
	Logout(c *fiber.Ctx, shaToken string) error
}

type authService struct {
	Log          *logrus.Logger
	DB           *gorm.DB
	Validate     *validator.Validate
	UserService  UserService
	TokenService TokenService
}

func NewAuthService(
	db *gorm.DB, userService UserService, tokenService TokenService,
) AuthService {
	return &authService{
		Log:          util.Log,
		DB:           db,
		Validate:     validation.Validator(),
		UserService:  userService,
		TokenService: tokenService,
	}
}

func (s *authService) Register(c *fiber.Ctx, req *request.Register) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	storeUser := &request.StoreUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}
	user, err := s.UserService.CreateUser(c, storeUser)
	if err != nil {
		s.Log.Errorf("Failed create user: %+v", err.Error())
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(c *fiber.Ctx, req *request.Login) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.UserService.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return user, nil
}

func (s *authService) Logout(c *fiber.Ctx, shaToken string) error {
	token, err := s.TokenService.GetTokenWithUser(c, shaToken)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Token not found")
	}
	userId := strconv.FormatUint(token.UserID, 10)
	err = s.TokenService.DeleteToken(c, userId)
	return err
}
