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
		createUser(c *fiber.Ctx) error
		updateUser(c *fiber.Ctx) error
		deleteUser(c *fiber.Ctx) error
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
	r.app.Post("/", r.createUser)
	r.app.Put("/:id<int>", r.updateUser)
	r.app.Delete("/:id<int>", r.deleteUser)
}

func (r *userRoute) getUsers(c *fiber.Ctx) error {
	res, err := r.userService.GetUsers()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取用户列表失败", err)
	}
	return domain.SuccessResponse(c, res, SuccessGetUsersMsg)
}

// 创建用户
func (r *userRoute) createUser(c *fiber.Ctx) error {
	params := new(domain.CreateUserParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := r.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := r.userService.CreateUser(*params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建用户失败", err)
	}

	return domain.SuccessResponse(c, nil, "创建用户成功")

}

// 删除用户
func (ur *userRoute) deleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数错误", err)
	}

	if err := ur.userService.DeleteUser(uint(id)); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除用户失败", err)
	}
	return domain.SuccessResponse(c, nil, "删除用户成功")
}

// 更新用户
func (ur *userRoute) updateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数错误", err)
	}
	params := new(domain.UpdateUserParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := ur.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	if err := ur.userService.UpdateUser(uint(id), *params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "更新用户失败", err)
	}
	return domain.SuccessResponse(c, nil, "更新用户成功")
}
