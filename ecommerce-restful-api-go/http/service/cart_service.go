package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"simple-api-go/config"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/validation"
	"simple-api-go/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CartService interface {
	GetCarts(c *fiber.Ctx) ([]response.CartResponse, error)
	AddCart(c *fiber.Ctx, req *request.CartStore) error
	RemoveCart(c *fiber.Ctx, productId string) error
	CartPayments(c *fiber.Ctx, req *request.CartPaymentStore) error
}

type cartService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Session  *session.Store
	Validate *validator.Validate
}

func NewCartService(db *gorm.DB) CartService {
	return &cartService{
		Log:      util.Log,
		DB:       db,
		Session:  session.New(),
		Validate: validation.Validator(),
	}
}

func (s *cartService) GetCarts(c *fiber.Ctx) ([]response.CartResponse, error) {
	carts := []response.CartResponse{}

	session, err := s.Session.Get(c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "session error")
	}

	authUser := util.AuthUser(c)
	sessionKeyUserId := fmt.Sprintf("%s_user_%v", config.SessionKey, authUser.ID)
	cart := session.Get(sessionKeyUserId)
	if cart == nil {
		return carts, nil
	}
	if err := json.Unmarshal([]byte(cart.(string)), &carts); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Cart data error")
	}
	return carts, nil
}

func (s *cartService) AddCart(c *fiber.Ctx, req *request.CartStore) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	authUser := util.AuthUser(c)
	quantity := req.Qty

	session, err := s.Session.Get(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Session data error")
	}

	sessionKeyUserId := fmt.Sprintf("%s_user_%v", config.SessionKey, authUser.ID)
	cartItems := []response.CartResponse{}
	cart := session.Get(sessionKeyUserId)
	if cart != nil {
		//DECODE DATA CART
		if err := json.Unmarshal([]byte(cart.(string)), &cartItems); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cart data error")
		}
	}

	// Check if cart exists
	if len(cartItems) == 0 && quantity <= 0 {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	// Check if product exists
	product := new(response.ProductResponse)
	result := s.DB.WithContext(c.Context()).First(product, "id = ?", req.ProductId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	// Check if item already exists in cart
	itemExists := false
	for i, item := range cartItems {
		if item.ProductId == product.Id {
			cartItems[i].Qty += quantity
			cartItems[i].TotalPrice = float64(cartItems[i].Qty) * product.Price
			itemExists = true
			if cartItems[i].Qty <= 0 {
				cartItems[i] = response.CartResponse{}
			}
			break
		}
	}

	if !itemExists {
		cartItems = append(cartItems, response.CartResponse{
			UserId:     authUser.ID,
			ProductId:  req.ProductId,
			Qty:        req.Qty,
			TotalPrice: float64(req.Qty) * product.Price,
			User: response.UserResponse{
				Id:    authUser.ID,
				Name:  authUser.Name,
				Email: authUser.Email,
			},
			Product: *product,
		})
	}

	//ENCODE DATA CART
	jsonData, errJson := json.Marshal(cartItems)
	if errJson != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	session.Set(sessionKeyUserId, string(jsonData))
	if err := session.Save(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to save session")
	}
	return err
}

func (s *cartService) RemoveCart(c *fiber.Ctx, productId string) error {
	session, err := s.Session.Get(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Session data error")
	}

	authUser := util.AuthUser(c)
	sessionKeyUserId := fmt.Sprintf("%s_user_%v", config.SessionKey, authUser.ID)

	cartItems := []response.CartResponse{}
	cart := session.Get(sessionKeyUserId)
	if cart != nil {
		//DECODE DATA CART
		if err := json.Unmarshal([]byte(cart.(string)), &cartItems); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cart data error")
		}
	}

	// Check if product exists
	product := new(response.ProductResponse)
	result := s.DB.WithContext(c.Context()).First(product, "id = ?", productId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	newCartItems := []response.CartResponse{}
	for _, item := range cartItems {
		if item.ProductId != product.Id && item.UserId == authUser.ID {
			newCartItems = append(newCartItems, item)
		}
	}

	//ENCODE DATA CART
	jsonData, errJson := json.Marshal(newCartItems)
	if errJson != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	session.Set(sessionKeyUserId, string(jsonData))
	if err := session.Save(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to save session")
	}

	return result.Error
}

func (s *cartService) CartPayments(c *fiber.Ctx, req *request.CartPaymentStore) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	authUser := util.AuthUser(c)
	sessionKeyUserId := fmt.Sprintf("%s_user_%v", config.SessionKey, authUser.ID)
	session, err := s.Session.Get(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Session data error")
	}

	cartItems := []response.CartResponse{}
	cart := session.Get(sessionKeyUserId)
	if cart != nil {
		//DECODE DATA CART
		if err := json.Unmarshal([]byte(cart.(string)), &cartItems); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cart data error")
		}
	}
	if req.PaymentMethod != "Tunai" && req.PaymentMethod != "QRIS" && req.PaymentMethod != "Transfer" {
		return fiber.NewError(fiber.StatusBadRequest, "Payment methods are not supported, only Cash, Transfer and QRIS")
	}

	tx := s.DB.Begin()
	newCartItems := []response.CartResponse{}

	//Create purchase data
	code := model.Purchase{}.GenerateCode("PRCH", tx)
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
		return fiber.NewError(fiber.StatusUnprocessableEntity, result.Error.Error())
	}

	purchaseDetails := []model.PurchaseDetail{}
	var total float64
	if len(req.CartDetails) > 0 {
		for _, v := range req.CartDetails {
			// Check if product exists
			product := new(response.ProductResponse)
			result := tx.WithContext(c.Context()).First(product, "id = ?", v.ProductId)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return fiber.NewError(fiber.StatusNotFound, "Product not found")
			}

			if v.Qty <= 0 {
				tx.Rollback()
				s.Log.Errorf("Failed to create purchase from cart qty not supported")
				return fiber.NewError(fiber.StatusBadRequest, "cart product qty not supported")
			}

			total += (float64(v.Qty) * product.Price)
			purchaseDetails = append(purchaseDetails, model.PurchaseDetail{
				PurchaseId: purchase.ID,
				ProductId:  product.Id,
				Qty:        v.Qty,
			})

			//update cart
			for _, item := range cartItems {
				if item.ProductId != v.ProductId && item.UserId == authUser.ID {
					newCartItems = append(newCartItems, item)
				}
			}
		}
	}
	//create purchase_details
	if res := tx.WithContext(c.Context()).Create(&purchaseDetails); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to create purchase details: %+v", res.Error)
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to create purchase details")

	}

	purchase.Total = total
	if res := tx.Save(&purchase); res.Error != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to update purchase: %+v", res.Error)
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to update purchase")
	}

	//ENCODE DATA CART
	jsonData, errJson := json.Marshal(newCartItems)
	if errJson != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	session.Set(sessionKeyUserId, string(jsonData))
	if err := session.Save(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to save session")
	}
	tx.Commit()
	return err
}
