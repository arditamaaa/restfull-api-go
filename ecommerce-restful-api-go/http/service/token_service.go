package service

import (
	"simple-api-go/config"
	"simple-api-go/http/model"
	"simple-api-go/http/validation"
	"simple-api-go/util"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TokenService interface {
	SaveToken(c *fiber.Ctx, token string, userID string, expires *time.Time) error
	DeleteToken(c *fiber.Ctx, userID string) error
	GetTokenWithUser(c *fiber.Ctx, tokenStr string) (*model.UserToken, error)
	GenerateAuthTokens(c *fiber.Ctx, user *model.User) (string, error)
}

type tokenService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewTokenService(db *gorm.DB) TokenService {
	return &tokenService{
		Log:      util.Log,
		DB:       db,
		Validate: validation.Validator(),
	}
}

func (s *tokenService) SaveToken(c *fiber.Ctx, token string, userID string, expires *time.Time) error {
	tx := s.DB.Begin()
	userId := util.StrToUint64(userID)

	userToken := model.UserToken{}
	tx.WithContext(c.Context()).
		Where("user_id = ?", userId).
		First(&userToken)

	userToken.Token = token
	userToken.Expires = expires
	userToken.UserID = userId
	if res := tx.Save(&userToken); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to update token: %+v", res.Error)
		return res.Error
	}
	tx.Commit()
	return nil
}

func (s *tokenService) DeleteToken(c *fiber.Ctx, userID string) error {
	tokenDoc := new(model.UserToken)
	result := s.DB.WithContext(c.Context()).
		Where("user_id = ?", userID).
		Delete(tokenDoc)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete token: %+v", result.Error)
	}
	return result.Error
}

func (s *tokenService) GetTokenWithUser(c *fiber.Ctx, tokenStr string) (*model.UserToken, error) {
	shaToken := util.EncodeSha256(tokenStr)
	tokenDoc := new(model.UserToken)
	result := s.DB.WithContext(c.Context()).
		Where("token = ?", shaToken).
		Preload("User").
		First(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed get token by user id: %+v", result.Error)
		return nil, result.Error
	}

	return tokenDoc, nil
}

func (s *tokenService) GenerateAuthTokens(c *fiber.Ctx, user *model.User) (string, error) {
	accessToken := util.RandStringBytes(40)
	shaToken := util.EncodeSha256(accessToken)
	var accessTokenExpires *time.Time
	if config.TokenExpMinutes > 0 {
		count := time.Now().UTC().Add(time.Minute * time.Duration(config.TokenExpMinutes))
		accessTokenExpires = &count

	}
	userID := strconv.FormatUint(user.ID, 10)
	if err := s.SaveToken(c, shaToken, userID, accessTokenExpires); err != nil {
		return "", err
	}
	return accessToken, nil
}
