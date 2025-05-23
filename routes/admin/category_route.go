package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	CategoryRoute interface {
		RegisterRoutes()
		getCategorys(c *fiber.Ctx) error
		createCategory(c *fiber.Ctx) error
		updateCategory(c *fiber.Ctx) error
		deleteCategory(c *fiber.Ctx) error
	}
	categoryRoute struct {
		app             fiber.Router
		categoryService services.CategoryService
		validator       *validator.Validate
	}
)

func NewCategoryRoute(app fiber.Router, categoryService services.CategoryService, validator *validator.Validate) CategoryRoute {
	return &categoryRoute{
		app,
		categoryService,
		validator,
	}
}

// 注册
func (r *categoryRoute) RegisterRoutes() {
	r.app.Get("/", r.getCategorys)
	r.app.Post("/", r.createCategory)
	r.app.Put("/:id<guid>", r.updateCategory)
	r.app.Delete("/:id<guid>", r.deleteCategory)
}

// 获取分类列表
func (r *categoryRoute) getCategorys(c *fiber.Ctx) error {
	res, err := r.categoryService.GetCategorys()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取分类列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取分类列表成功")
}

// 创建分类
func (r *categoryRoute) createCategory(c *fiber.Ctx) error {
	params := new(domain.CreateCategoryParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.categoryService.CreateCategory(*params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建分类失败", err)
	}
	return domain.SuccessResponse(c, nil, "创建分类成功")
}

// 更新分类
func (r *categoryRoute) updateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

	params := new(domain.UpdateCategoryParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.categoryService.UpdateCategory(id, *params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "更新分类失败", err)
	}
	return domain.SuccessResponse(c, nil, "更新分类成功")
}

// 删除分类
func (r *categoryRoute) deleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}
	if err := r.categoryService.DeleteCategory(id); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除分类失败", err)
	}
	return domain.SuccessResponse(c, nil, "删除分类成功")
}
