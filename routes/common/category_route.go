package common

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
		getCategoryByAlias(c *fiber.Ctx) error
		getCategoryByID(c *fiber.Ctx) error
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
	r.app.Get("/getCategoryByAlias/:alias", r.getCategoryByAlias)
	r.app.Get("/getCategoryByID/:id", r.getCategoryByID)
}

// 获取分类列表
func (r *categoryRoute) getCategorys(c *fiber.Ctx) error {
	res, err := r.categoryService.GetCategorysWithCache()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取分类列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取分类列表成功")
}

// 根据别名获取分类
func (r *categoryRoute) getCategoryByAlias(c *fiber.Ctx) error {
	alias := c.Params("alias")
	category, err := r.categoryService.GetCategoryByAliasWithCache(alias)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取分类失败", err)
	}
	return domain.SuccessResponse(c, category, "获取分类成功")
}

// 根据文章ID获取分类
func (r *categoryRoute) getCategoryByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "无效的文章ID", err)
	}
	category, err := r.categoryService.GetCategoryByIDWithCache(id)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取分类失败", err)
	}
	return domain.SuccessResponse(c, category, "获取分类成功")
}
