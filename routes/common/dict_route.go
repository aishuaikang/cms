package common

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	DictRoute interface {
		RegisterRoutes()
		getDictExtraByCode(c *fiber.Ctx) error
		getDictByCode(c *fiber.Ctx) error
		getSubDictsByCode(c *fiber.Ctx) error
	}
	dictRoute struct {
		app         fiber.Router
		dictService services.DictService
		validator   *validator.Validate
	}
)

func NewDictRoute(app fiber.Router, dictService services.DictService, validator *validator.Validate) DictRoute {
	return &dictRoute{
		app,
		dictService,
		validator,
	}
}

// 注册
func (r *dictRoute) RegisterRoutes() {
	r.app.Get("/getDictExtraByCode/:code", r.getDictExtraByCode)
	r.app.Get("/getDictByCode/:code", r.getDictByCode)
	r.app.Get("/getSubDictsByCode/:code", r.getSubDictsByCode)
}

// 根据code获取字典的extra
func (r *dictRoute) getDictExtraByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	extra, err := r.dictService.GetDictExtraByCodeWithCache(code)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取字典extra失败", err)
	}
	return domain.SuccessResponse(c, extra, "获取字典extra成功")
}

// 根据code获取字典的extra
func (r *dictRoute) getDictByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	extra, err := r.dictService.GetDictByCodeWithCache(code)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取字典失败", err)
	}
	return domain.SuccessResponse(c, extra, "获取字典成功")
}

// 根据code获取子字典列表
func (r *dictRoute) getSubDictsByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	subDicts, err := r.dictService.GetSubDictsByCodeWithCache(code)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取子字典列表失败", err)
	}
	return domain.SuccessResponse(c, subDicts, "获取子字典列表成功")
}
