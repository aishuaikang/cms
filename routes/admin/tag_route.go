package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	TagRoute interface {
		RegisterRoutes()
		getTags(c *fiber.Ctx) error
		createTag(c *fiber.Ctx) error
		updateTag(c *fiber.Ctx) error
		deleteTag(c *fiber.Ctx) error
	}
	tagRoute struct {
		app        fiber.Router
		tagService services.TagService
		validator  *validator.Validate
	}
)

func NewTagRoute(app fiber.Router, tagService services.TagService, validator *validator.Validate) TagRoute {
	return &tagRoute{
		app,
		tagService,
		validator,
	}
}

// 注册
func (r *tagRoute) RegisterRoutes() {
	r.app.Get("/", r.getTags)
	r.app.Post("/", r.createTag)
	r.app.Put("/:id<guid>", r.updateTag)
	r.app.Delete("/:id<guid>", r.deleteTag)
}

// 获取标签列表
func (r *tagRoute) getTags(c *fiber.Ctx) error {
	res, err := r.tagService.GetTags()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取标签列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取标签列表成功")
}

// 创建标签
func (r *tagRoute) createTag(c *fiber.Ctx) error {
	params := new(domain.CreateTagParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.tagService.CreateTag(*params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建标签失败", err)
	}

	return domain.SuccessResponse(c, nil, "创建标签成功")
}

// 更新标签
func (r *tagRoute) updateTag(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}
	params := new(domain.UpdateTagParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.tagService.UpdateTag(id, *params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "更新标签失败", err)
	}

	return domain.SuccessResponse(c, nil, "更新标签成功")
}

// 删除标签
func (r *tagRoute) deleteTag(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

	if err := r.tagService.DeleteTag(id); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除标签失败", err)
	}

	return domain.SuccessResponse(c, nil, "删除标签成功")
}
