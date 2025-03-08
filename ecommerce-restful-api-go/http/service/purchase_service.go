package service

import (
	"errors"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/validation"
	"simple-api-go/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PurchaseService interface {
	GetPurchases(c *fiber.Ctx, params *request.QueryPurchase) ([]model.Purchase, int64, error)
	GetPurchaseByID(c *fiber.Ctx, id string) (*model.Purchase, error)
	CreatePurchase(c *fiber.Ctx, req *request.PurchaseStore) (*model.Purchase, error)
	UpdatePurchase(c *fiber.Ctx, req *request.PurchaseStore, id string) (*model.Purchase, error)
	DeletePurchase(c *fiber.Ctx, id string) error
}

type purchaseService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewPurchaseService(db *gorm.DB) PurchaseService {
	return &purchaseService{
		Log:      util.Log,
		DB:       db,
		Validate: validation.Validator(),
	}
}

func (s *purchaseService) GetPurchases(c *fiber.Ctx, params *request.QueryPurchase) ([]model.Purchase, int64, error) {
	var purchases []model.Purchase
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("code LIKE ?", "%"+search+"%")
	}
	if userId := params.UserId; userId != 0 {
		query = query.Where("user_id = ?", userId)
	}

	result := query.Find(&purchases).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search purchases: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Limit(params.Limit).Offset(offset).Preload("User").Find(&purchases)
	if result.Error != nil {
		s.Log.Errorf("Failed to get all purchases: %+v", result.Error)
		return nil, 0, result.Error
	}

	return purchases, totalResults, result.Error
}

func (s *purchaseService) GetPurchaseByID(c *fiber.Ctx, id string) (*model.Purchase, error) {
	purchase := new(model.Purchase)

	result := s.DB.WithContext(c.Context()).
		Preload("PurchaseDetails").
		Preload("PurchaseDetails.Product").
		First(purchase, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "purchase not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed get purchase by id: %+v", result.Error)
	}
	return purchase, result.Error
}

func (s *purchaseService) CreatePurchase(c *fiber.Ctx, req *request.PurchaseStore) (*model.Purchase, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	authUser := util.AuthUser(c)
	tx := s.DB.Begin()

	if req.PaymentMethod != "Tunai" &&
		req.PaymentMethod != "QRIS" &&
		req.PaymentMethod != "Transfer" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Payment methods are not supported, only Cash, Transfer and QRIS")
	}

	//Create purchase data
	code := model.Purchase{}.GenerateCode("PO", tx)
	purchase := &model.Purchase{
		Code:          code,
		UserID:        authUser.ID,
		PaymentMethod: req.PaymentMethod,
		Status:        "DONE",
	}

	result := tx.WithContext(c.Context()).Create(purchase)
	if result.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to create purchase: %+v", result.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, result.Error.Error())
	}

	purchaseDetails := []model.PurchaseDetail{}
	var total float64
	if len(req.PurchaseDetails) > 0 {
		for _, v := range req.PurchaseDetails {
			// Check if product exists
			product := new(response.ProductResponse)
			result := tx.WithContext(c.Context()).First(product, "id = ?", v.ProductId)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
			}

			if v.Qty <= 0 {
				tx.Rollback()
				s.Log.Errorf("Failed to create purchase from cart qty not supported")
				return nil, fiber.NewError(fiber.StatusBadRequest, "cart product qty not supported")
			}

			total += (float64(v.Qty) * product.Price)
			purchaseDetails = append(purchaseDetails, model.PurchaseDetail{
				PurchaseId: purchase.ID,
				ProductId:  product.Id,
				Qty:        v.Qty,
			})
		}
	}

	//create purchase_details
	if res := tx.WithContext(c.Context()).Create(&purchaseDetails); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to create purchase details: %+v", res.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to create purchase details")

	}

	purchase.Total = total
	if res := tx.Save(&purchase); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to update purchase: %+v", res.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to update purchase")
	}
	tx.Commit()
	return purchase, result.Error
}

func (s *purchaseService) UpdatePurchase(c *fiber.Ctx, req *request.PurchaseStore, id string) (*model.Purchase, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	tx := s.DB.Begin()
	purchase := new(model.Purchase)
	result := tx.WithContext(c.Context()).First(purchase, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "purchase not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed get purchase by id: %+v", result.Error)
	}

	if req.PaymentMethod != "Tunai" &&
		req.PaymentMethod != "QRIS" &&
		req.PaymentMethod != "Transfer" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Payment methods are not supported, only Cash, Transfer and QRIS")
	}
	purchase.PaymentMethod = req.PaymentMethod

	if res := tx.Delete(&model.PurchaseDetail{}, "purchase_id =?", purchase.ID); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to delete purchase details: %+v", res.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to de;ete purchase details")

	}

	purchaseDetails := []model.PurchaseDetail{}
	var total float64
	if len(req.PurchaseDetails) > 0 {
		for _, v := range req.PurchaseDetails {
			// Check if product exists
			product := new(response.ProductResponse)
			result := tx.WithContext(c.Context()).First(product, "id = ?", v.ProductId)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
			}

			if v.Qty <= 0 {
				tx.Rollback()
				s.Log.Errorf("Failed to create purchase from cart qty not supported")
				return nil, fiber.NewError(fiber.StatusBadRequest, "cart product qty not supported")
			}

			total += (float64(v.Qty) * product.Price)
			purchaseDetails = append(purchaseDetails, model.PurchaseDetail{
				PurchaseId: purchase.ID,
				ProductId:  product.Id,
				Qty:        v.Qty,
			})
		}
	}

	//create purchase_details
	if res := tx.WithContext(c.Context()).Create(&purchaseDetails); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to create purchase details: %+v", res.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to create purchase details")

	}

	purchase.Total = total
	if res := tx.Save(&purchase); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to update purchase: %+v", res.Error)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to update purchase")
	}

	tx.Commit()

	return purchase, result.Error
}

func (s *purchaseService) DeletePurchase(c *fiber.Ctx, id string) error {
	purchase := new(model.Purchase)

	tx := s.DB.Begin()
	result := tx.WithContext(c.Context()).First(purchase, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "purchase not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed get purchase by id: %+v", result.Error)
	}

	purchase.Status = "CANCEL"
	if res := tx.Save(&purchase); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to update purchase: %+v", res.Error)
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to update purchase")
	}
	if res := tx.Delete(&purchase); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to delete purchase: %+v", res.Error)
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to delete purchase")
	}
	tx.Commit()
	return result.Error
}
