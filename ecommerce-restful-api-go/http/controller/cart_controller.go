package controller

import (
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

type CartController struct {
	CartService service.CartService
}

func NewCartController(cartService service.CartService) *CartController {
	return &CartController{
		CartService: cartService,
	}
}

func (u *CartController) Index(c *fiber.Ctx) error {
	carts, err := u.CartService.GetCarts(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).
		JSON(response.SuccessList[response.CartResponse]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get all cart successfully",
			Data:    carts,
		})
}

func (u *CartController) CreateOrUpdate(c *fiber.Ctx) error {
	req := new(request.CartStore)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := u.CartService.AddCart(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Product added to cart",
		})
}

func (u *CartController) Delete(c *fiber.Ctx) error {
	productID := c.Params("productId")
	if err := u.CartService.RemoveCart(c, productID); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Remove product from cart successfully",
		})
}

func (u *CartController) Pay(c *fiber.Ctx) error {
	req := new(request.CartPaymentStore)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := u.CartService.CartPayments(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Payment successful",
		})
}
