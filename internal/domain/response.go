package domain

import (
	"github.com/gofiber/fiber"
)

type WebResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func NewSuccessfulResponse(c *fiber.Ctx, data interface{}) error {
	c.SendStatus(fiber.StatusOK)
	return c.JSON(WebResponse{
		Message: "Sukses",
		Code:    fiber.StatusOK,
		Data:    data,
	})
}

func NewBadRequestResponse(c *fiber.Ctx) error {
	c.SendStatus(fiber.StatusBadRequest)
	return c.JSON(WebResponse{
		Message: "Bad Request",
		Code:    fiber.StatusBadRequest,
		Data:    []string{},
	})
}

func NewForbiddenResponse(c *fiber.Ctx) error {
	c.SendStatus(fiber.StatusForbidden)
	return c.JSON(WebResponse{
		Message: "Bad Request",
		Code:    fiber.StatusForbidden,
		Data:    []string{},
	})
}

func NewUnexpectedErrorRequest(c *fiber.Ctx, msg string) error {
	c.SendStatus(fiber.StatusForbidden)
	return c.JSON(WebResponse{
		Message: msg,
		Code:    fiber.StatusForbidden,
		Data:    []string{},
	})
}

func NewDefaultErrorResponse(c *fiber.Ctx, msg string, code int) error {
	c.SendStatus(fiber.StatusBadRequest)
	return c.JSON(WebResponse{
		Message: msg,
		Code:    code,
		Data:    []string{},
	})
}

func NewUnexpectedErrorResponse(c *fiber.Ctx, m string) error {
	c.SendStatus(fiber.StatusInternalServerError)

	return c.JSON(WebResponse{
		Message: m,
		Code:    fiber.StatusInternalServerError,
		Data:    []string{},
	})
}

func NewUnauthorizedResponse(c *fiber.Ctx) error {
	c.SendStatus(fiber.StatusUnauthorized)

	return c.JSON(WebResponse{
		Message: "Unauthorized",
		Code:    fiber.StatusUnauthorized,
		Data:    []string{},
	})
}

func NewDefaultBrigateResponse(c *fiber.Ctx, m string, code int) error {
	c.SendStatus(fiber.StatusOK)

	return c.JSON(WebResponse{
		Message: m,
		Code:    code,
		Data:    []string{},
	})
}

type ResponseError struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
}
