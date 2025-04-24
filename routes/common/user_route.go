package common

import (
	"cms/models/domain"
	"cms/services"
	"crypto/rsa"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type (
	AccountRoute interface {
		RegisterRoutes()
		login(c *fiber.Ctx) error
	}
	accountRoute struct {
		app         fiber.Router
		validator   *validator.Validate
		userService services.UserService
		privateKey  *rsa.PrivateKey
	}
)

func NewAccountRoute(app fiber.Router, userService services.UserService, validator *validator.Validate, privateKey *rsa.PrivateKey) AccountRoute {
	return &accountRoute{
		app:         app,
		validator:   validator,
		userService: userService,
		privateKey:  privateKey,
	}
}
func (ur *accountRoute) RegisterRoutes() {
	ur.app.Post("/login", ur.login)
}

func (ur *accountRoute) login(c *fiber.Ctx) error {
	params := new(domain.LoginParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := ur.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	// // Throws Unauthorized error
	// if params.Username != "admin" || params.Password != "123456" {
	// 	return domain.ErrorResponse(c, fiber.StatusUnauthorized, "用户名或密码错误", ErrInvalidCredentials)
	// }

	user, err := ur.userService.Login(*params)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusUnauthorized, "用户名或密码错误", err)
	}

	// 设置 JWT token 的过期时间
	const tokenExpiration = time.Hour * 1

	// 创建 JWT token 的声明
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
	}

	// 创建一个新的 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// 使用私钥签名 token
	// 这里使用了 RS512 签名算法
	t, err := token.SignedString(ur.privateKey)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "生成 JWT 失败", err)
	}

	res := domain.LoginResponse{
		User:  *user,
		Token: t,
	}

	return domain.SuccessResponse(c, res, "登录成功")

}
