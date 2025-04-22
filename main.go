package main

import (
	"cms/config"
	"cms/models/domain"
	"cms/routes/admin"
	"cms/routes/common"
	"cms/services"
	"cms/utils"

	"github.com/bytedance/sonic"

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

	systemConfig, err := config.NewSystemConfig()
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
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	// app.Use(csrf.New())

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := app.Group("api")

	userService := services.NewUserService(db)

	// 创建初始用户
	if err := userService.CreateInitialUser(systemConfig.SysAdminUser, systemConfig.SysAdminPassword); err != nil {
		panic(err)
	}

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
		admin.NewUserRoute(adminGroup.Group("user"), userService, validate).RegisterRoutes()
		// 标签
		admin.NewTagRoute(adminGroup.Group("tag"), services.NewTagService(db), validate).RegisterRoutes()
		// 字典
		admin.NewDictRoute(adminGroup.Group("dict"), services.NewDictService(db), validate).RegisterRoutes()
	}

	{
		commonGroup := api.Group("common")
		// 图片
		common.NewImageRoute(commonGroup.Group("image"), services.NewImageService(db), validate).RegisterRoutes()
		// 用户
		common.NewUserRoute(commonGroup.Group("user"), userService, validate, privateKey).RegisterRoutes()
	}

	app.Listen(":3000")
}
