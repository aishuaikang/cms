package admin

import (
	"cms/models/domain"
	"cms/services"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// 获取用户ID失败
	ErrGetUserIDFailed = errors.New("获取用户ID失败")
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
	r.app.Put("/:id<guid>", r.updateArticle)
	r.app.Delete("/:id<guid>", r.deleteArticle)
}

// 获取文章列表
func (r *articleRoute) getArticles(c *fiber.Ctx) error {
	params := new(domain.GetArticleListParams)
	if err := c.QueryParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析查询参数失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	res, err := r.articleService.GetArticles(*params)
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

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "获取用户ID失败", ErrGetUserIDFailed)
	}

	if err := r.articleService.CreateArticle(userID, *params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建文章失败", err)
	}
	return domain.SuccessResponse(c, nil, "创建文章成功")
}

// 更新文章
func (r *articleRoute) updateArticle(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

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
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析ID失败", err)
	}

	if err := r.articleService.DeleteArticle(id); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除文章失败", err)
	}
	return domain.SuccessResponse(c, nil, "删除文章成功")
}
