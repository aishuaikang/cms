package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

type (
	AccountRoute interface {
		RegisterRoutes()
		logout(c *fiber.Ctx) error
	}
	accountRoute struct {
		app         fiber.Router
		validator   *validator.Validate
		userService services.UserService
	}
)

func NewAccountRoute(app fiber.Router, userService services.UserService, validator *validator.Validate) AccountRoute {
	return &accountRoute{
		app:         app,
		validator:   validator,
		userService: userService,
	}
}

// RegisterRoutes 注册路由
func (ar *accountRoute) RegisterRoutes() {
	ar.app.Post("/logout", ar.logout)
}

// 登出
func (ar *accountRoute) logout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	log.Info("用户 %s 登出", user)
	return domain.SuccessResponse(c, fiber.Map{}, "登出成功")
}
