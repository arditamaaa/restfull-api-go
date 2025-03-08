package controller

import (
	"math"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (u *ProductController) Index(c *fiber.Ctx) error {
	query := &request.QueryProduct{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("paginate", 10),
		Search: c.Query("keyword", ""),
	}

	products, totalResults, err := u.ProductService.GetProducts(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.Product]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all products successfully",
			Results:      products,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

func (s *ProductController) Show(c *fiber.Ctx) error {
	productID := c.Params("id")
	product, err := s.ProductService.GetProductByID(c, productID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.Product]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get Product successfully",
			Data:    *product,
		})
}

func (s *ProductController) Create(c *fiber.Ctx) error {
	req := new(request.ProductStore)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	product, err := s.ProductService.CreateProduct(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.Success[model.Product]{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create Product successfully",
			Data:    *product,
		})
}

func (s *ProductController) Update(c *fiber.Ctx) error {
	req := new(request.ProductStore)
	productID := c.Params("id")

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	product, err := s.ProductService.UpdateProduct(c, req, productID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.Product]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update product successfully",
			Data:    *product,
		})
}

func (s *ProductController) Delete(c *fiber.Ctx) error {
	productID := c.Params("id")
	if err := s.ProductService.DeleteProduct(c, productID); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete product successfully",
		})
}
