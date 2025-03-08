package service

import (
	"errors"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/validation"
	"simple-api-go/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService interface {
	GetProducts(c *fiber.Ctx, params *request.QueryProduct) ([]model.Product, int64, error)
	GetProductByID(c *fiber.Ctx, id string) (*model.Product, error)
	CreateProduct(c *fiber.Ctx, req *request.ProductStore) (*model.Product, error)
	UpdateProduct(c *fiber.Ctx, req *request.ProductStore, id string) (*model.Product, error)
	DeleteProduct(c *fiber.Ctx, id string) error
}

type productService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewProductService(db *gorm.DB) ProductService {
	return &productService{
		Log:      util.Log,
		DB:       db,
		Validate: validation.Validator(),
	}
}

func (s *productService) GetProducts(c *fiber.Ctx, params *request.QueryProduct) ([]model.Product, int64, error) {
	var products []model.Product
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	result := query.Find(&products).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search products: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Limit(params.Limit).Offset(offset).Find(&products)
	if result.Error != nil {
		s.Log.Errorf("Failed to get all products: %+v", result.Error)
		return nil, 0, result.Error
	}

	return products, totalResults, result.Error
}

func (s *productService) GetProductByID(c *fiber.Ctx, id string) (*model.Product, error) {
	product := new(model.Product)

	result := s.DB.WithContext(c.Context()).First(product, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get user by id: %+v", result.Error)
	}

	return product, result.Error
}

func (s *productService) CreateProduct(c *fiber.Ctx, req *request.ProductStore) (*model.Product, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	product := &model.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	result := s.DB.WithContext(c.Context()).Create(product)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Product is already exist")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to create product: %+v", result.Error)
	}

	return product, result.Error
}

func (s *productService) UpdateProduct(c *fiber.Ctx, req *request.ProductStore, id string) (*model.Product, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	product := &model.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	result := s.DB.WithContext(c.Context()).Where("id = ?", id).Updates(product)
	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to update product: %+v", result.Error)
	}
	return product, result.Error
}

func (s *productService) DeleteProduct(c *fiber.Ctx, id string) error {
	product := new(model.Product)

	result := s.DB.WithContext(c.Context()).Delete(product, "id = ?", id)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to delete product: %+v", result.Error)
	}

	return result.Error
}
