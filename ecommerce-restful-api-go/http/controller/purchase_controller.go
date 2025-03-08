package controller

import (
	"math"
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

type PurchaseController struct {
	Service service.PurchaseService
}

func NewPurchaseController(purchaseService service.PurchaseService) *PurchaseController {
	return &PurchaseController{
		Service: purchaseService,
	}
}

func (u *PurchaseController) Index(c *fiber.Ctx) error {
	query := &request.QueryPurchase{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("paginate", 10),
		Search: c.Query("keyword", ""),
		UserId: c.QueryInt("user_id"),
	}

	purchases, totalResults, err := u.Service.GetPurchases(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.Purchase]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all purchase successfully",
			Results:      purchases,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

func (s *PurchaseController) Show(c *fiber.Ctx) error {
	purchaseID := c.Params("id")
	purchase, err := s.Service.GetPurchaseByID(c, purchaseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.Purchase]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get Product successfully",
			Data:    *purchase,
		})
}

func (s *PurchaseController) Create(c *fiber.Ctx) error {
	req := new(request.PurchaseStore)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	purchase, err := s.Service.CreatePurchase(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.Success[model.Purchase]{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create purchase successfully",
			Data:    *purchase,
		})
}

func (s *PurchaseController) Update(c *fiber.Ctx) error {
	req := new(request.PurchaseStore)
	purchaseID := c.Params("id")

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	product, err := s.Service.UpdatePurchase(c, req, purchaseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success[model.Purchase]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update purchase successfully",
			Data:    *product,
		})
}

func (s *PurchaseController) Delete(c *fiber.Ctx) error {
	purchaseID := c.Params("id")
	if err := s.Service.DeletePurchase(c, purchaseID); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete purchase successfully",
		})
}
