package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	DictRoute interface {
		RegisterRoutes()
		getDicts(c *fiber.Ctx) error
		createDict(c *fiber.Ctx) error
		updateDict(c *fiber.Ctx) error
		deleteDict(c *fiber.Ctx) error
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
	r.app.Get("/", r.getDicts)
	r.app.Post("/", r.createDict)
	r.app.Put("/:id<guid>", r.updateDict)
	r.app.Delete("/:id<guid>", r.deleteDict)
}

// 获取字典列表
func (r *dictRoute) getDicts(c *fiber.Ctx) error {
	res, err := r.dictService.GetDicts()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取字典列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取字典列表成功")
}

// 创建字典
func (r *dictRoute) createDict(c *fiber.Ctx) error {
	params := new(domain.CreateDictParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.dictService.CreateDict(*params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建字典失败", err)
	}
	return domain.SuccessResponse(c, nil, "创建字典成功")
}

// 更新字典
func (r *dictRoute) updateDict(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

	params := new(domain.UpdateDictParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.dictService.UpdateDict(id, *params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "更新字典失败", err)
	}
	return domain.SuccessResponse(c, nil, "更新字典成功")
}

// 删除字典
func (r *dictRoute) deleteDict(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

	if err := r.dictService.DeleteDict(id); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除字典失败", err)
	}
	return domain.SuccessResponse(c, nil, "删除字典成功")
}
