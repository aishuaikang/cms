package admin

import (
	"cms/models/domain"
	"cms/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	// 获取用户列表成功
	SuccessGetUsersMsg = "获取用户列表成功"
)

type (
	UserRoute interface {
		RegisterRoutes()
		getUsers(c *fiber.Ctx) error
	}
	userRoute struct {
		app         fiber.Router
		validator   *validator.Validate
		userService services.UserService
	}
)

func NewUserRoute(app fiber.Router, userService services.UserService, validator *validator.Validate) UserRoute {
	return &userRoute{
		app:         app,
		validator:   validator,
		userService: userService,
	}
}

// RegisterRoutes 注册路由
func (r *userRoute) RegisterRoutes() {
	r.app.Get("/", r.getUsers)
	// r.app.Post("/", r.createCategory)
	// r.app.Put("/:id", r.updateCategory)
	// r.app.Delete("/:id", r.deleteCategory)
}

func (r *userRoute) getUsers(c *fiber.Ctx) error {
	res, err := r.userService.GetUsers()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取用户列表失败", err)
	}
	return domain.SuccessResponse(c, res, SuccessGetUsersMsg)
}
