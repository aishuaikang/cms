package main

import (
	"cms/models/domain"
	"cms/routes/admin"
	"cms/routes/common"
	"cms/services"
	"cms/utils"

	"github.com/goccy/go-json"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := utils.InitDB()
	if err != nil {
		panic(err)
	}

	privateKey, err := utils.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	validate := validator.New()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "cms",
		AppName:      "cms v0.0.1",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := app.Group("api")

	{
		// 对于所有admin路由，使用jwt中间件进行验证
		// 这里的jwt中间件会在请求到达路由之前进行验证
		adminGroup := api.Group("admin", jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: privateKey.Public()},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return domain.ErrorResponse(c, fiber.StatusUnauthorized, "未授权", err)
			},
		}))

		// 分类
		admin.NewCategoryRoute(adminGroup.Group("category"), services.NewCategoryService(db), validate).RegisterRoutes()
		// 文章
		admin.NewArticleRoute(adminGroup.Group("article"), services.NewArticleService(db), validate).RegisterRoutes()
		// 图片
		admin.NewImageRoute(adminGroup.Group("image"), services.NewImageService(db), validate).RegisterRoutes()
		// 用户
		admin.NewUserRoute(adminGroup.Group("user"), services.NewUserService(db), validate).RegisterRoutes()
	}

	{
		commonGroup := api.Group("common")
		// 图片
		common.NewImageRoute(commonGroup.Group("image"), services.NewImageService(db), validate).RegisterRoutes()
		// 用户
		common.NewUserRoute(commonGroup.Group("user"), validate, privateKey).RegisterRoutes()
	}

	app.Listen(":3000")
}
