package domain

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   any    `json:"error"`
	Data    any    `json:"data"`
}

type LimitResponse[T any] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
	Pages int   `json:"pages"`
}

func SuccessResponse(ctx *fiber.Ctx, data any, message string) error {
	resp := Response{
		Code:    200,
		Message: message,
		Data:    data,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func ErrorResponse(ctx *fiber.Ctx, code int, message string, err error) error {
	resp := Response{
		Code:    code,
		Message: message,
		Error:   err.Error(),
		Data:    nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
