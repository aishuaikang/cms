package roleauth

import (
	"cms/models/domain"
	"cms/services"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("用户未找到")
)

func New(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return domain.ErrorResponse(c, fiber.StatusBadRequest, "获取用户ID失败", ErrUserNotFound)
		}

		isSuper, err := userService.GetUserIsSuper(userID)
		if err != nil {
			return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取用户角色失败", err)
		}

		if !isSuper {
			return domain.ErrorResponse(c, fiber.StatusForbidden, "没有权限访问该资源", errors.New("没有权限访问该资源"))
		}

		return c.Next()
	}
}
