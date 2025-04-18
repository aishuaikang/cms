package common

import (
	"cms/models/domain"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidCredentials 用户名或密码错误
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

type (
	UserRoute interface {
		RegisterRoutes()
		login(c *fiber.Ctx) error
	}
	userRoute struct {
		app        fiber.Router
		validator  *validator.Validate
		privateKey *rsa.PrivateKey
	}
)

func NewUserRoute(app fiber.Router, validator *validator.Validate, privateKey *rsa.PrivateKey) UserRoute {
	return &userRoute{
		app:        app,
		validator:  validator,
		privateKey: privateKey,
	}
}
func (ur *userRoute) RegisterRoutes() {
	ur.app.Post("/login", ur.login)
}

func (ur *userRoute) login(c *fiber.Ctx) error {
	params := new(domain.LoginParams)
	if err := c.BodyParser(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "解析请求体失败", err)
	}

	if err := ur.validator.Struct(params); err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数校验失败", err)
	}

	// Throws Unauthorized error
	if params.Username != "admin" || params.Password != "123456" {
		return domain.ErrorResponse(c, fiber.StatusUnauthorized, "用户名或密码错误", ErrInvalidCredentials)
	}

	// Create the Claims
	const tokenExpiration = time.Second * 60
	claims := jwt.MapClaims{
		"name": params.Username,
		"exp":  time.Now().Add(tokenExpiration).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(ur.privateKey)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "生成 JWT 失败", err)
	}
	return domain.SuccessResponse(c, fiber.Map{
		"token": t,
	}, "登录成功")

}
