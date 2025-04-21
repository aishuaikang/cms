package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ArticleRoute interface {
		RegisterRoutes()
		getArticles(c *fiber.Ctx) error
		createArticle(c *fiber.Ctx) error
		updateArticle(c *fiber.Ctx) error
		deleteArticle(c *fiber.Ctx) error
	}
	articleRoute struct {
		app            fiber.Router
		articleService services.ArticleService
		validator      *validator.Validate
	}
)

func NewArticleRoute(app fiber.Router, articleService services.ArticleService, validator *validator.Validate) ArticleRoute {
	return &articleRoute{
		app,
		articleService,
		validator,
	}
}

// 注册
func (r *articleRoute) RegisterRoutes() {
	r.app.Get("/", r.getArticles)
	r.app.Post("/", r.createArticle)
	r.app.Put("/:id", r.updateArticle)
	r.app.Delete("/:id", r.deleteArticle)
}

// 获取文章列表
func (r *articleRoute) getArticles(c *fiber.Ctx) error {
	res, err := r.articleService.GetArticles()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取文章列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取文章列表成功")
}

// 创建文章
func (r *articleRoute) createArticle(c *fiber.Ctx) error {
	params := new(domain.CreateArticleParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.articleService.CreateArticle(*params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建文章失败", err)
	}
	return domain.SuccessResponse(c, nil, "创建文章成功")
}

// 更新文章
func (r *articleRoute) updateArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	body := new(domain.UpdateArticleParams)
	if err := c.BodyParser(body); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(body); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.articleService.UpdateArticle(id, *body); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "更新文章失败", err)
	}

	return domain.SuccessResponse(c, nil, "更新文章成功")
}

// 删除文章
func (r *articleRoute) deleteArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	err := r.articleService.DeleteArticle(id)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除文章失败", err)
	}
	return domain.SuccessResponse(c, nil, "删除文章成功")
}
