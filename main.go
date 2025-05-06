package main

import (
	"cms/config"
	"cms/middleware/roleauth"
	"cms/models/domain"
	"cms/models/scopes"
	"cms/routes/admin"
	"cms/routes/common"
	"cms/services"
	"cms/utils"
	"os"
	"time"

	"github.com/bytedance/sonic"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
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
		// Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "cms",
		AppName:      "cms v0.0.1",
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		// ErrorHandler: func(c *fiber.Ctx, err error) error {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
		// 		Code:    fiber.StatusInternalServerError,
		// 		Message: "服务器错误",
		// 		Data:    nil,
		// 		Error:   err.Error(),
		// 	})
		// },
	})

	// 设置压缩中间件
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 2
	}))

	// 设置头盔中间件
	// 头盔中间件可以帮助你设置一些HTTP头部来增强安全性
	// 例如，设置X-Content-Type-Options、X-Frame-Options等
	app.Use(helmet.New())

	// 设置日志中间件
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/", monitor.New(monitor.Config{
		Title: "CMS 指标监控",
	}))

	api := app.Group("api")

	userService := services.NewUserService(db)

	// 创建初始用户
	if err := userService.CreateInitialUser(systemConfig.SysAdminUser, systemConfig.SysAdminPassword); err != nil {
		panic(err)
	}

	roleAuthMiddleware := roleauth.New(userService)

	categoryService := services.NewCategoryService(db)
	articleService := services.NewArticleService(db, scopes.NewArticleScope(db))
	imageService := services.NewImageService(db)
	dictService := services.NewDictService(db)

	{
		// 对于所有admin路由，使用jwt中间件进行验证
		// 这里的jwt中间件会在请求到达路由之前进行验证
		adminGroup := api.Group("admin", jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{
				JWTAlg: jwtware.RS512,
				Key:    privateKey.Public()},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return domain.ErrorResponse(c, fiber.StatusUnauthorized, "您的身份验证已过期，请重新登录", err)
			},
		}))

		// 分类
		admin.NewCategoryRoute(adminGroup.Group("category", roleAuthMiddleware), categoryService, validate).RegisterRoutes()
		// 文章
		admin.NewArticleRoute(adminGroup.Group("article"), articleService, validate).RegisterRoutes()
		// 图片
		admin.NewImageRoute(adminGroup.Group("image"), imageService, validate).RegisterRoutes()
		// 用户
		admin.NewUserRoute(adminGroup.Group("user", roleAuthMiddleware), userService, validate).RegisterRoutes()
		// 标签
		admin.NewTagRoute(adminGroup.Group("tag"), services.NewTagService(db), validate).RegisterRoutes()
		// 字典
		admin.NewDictRoute(adminGroup.Group("dict", roleAuthMiddleware), dictService, validate).RegisterRoutes()
		// 账号
		admin.NewAccountRoute(adminGroup.Group("account"), userService, validate).RegisterRoutes()
	}

	{
		// 设置限制器中间件
		// 限制器中间件可以帮助你限制请求的频率
		// 例如，限制每个IP每分钟只能请求100次
		// 这里使用了滑动窗口算法

		commonGroup := api.Group("common", limiter.New(limiter.Config{
			Max:        1000,            // 每个IP的最大请求数
			Expiration: 1 * time.Minute, // 时间窗口为1分钟
			// KeyGenerator: func(c *fiber.Ctx) string {
			// 	return c.IP() // 使用IP地址作为限流的唯一标识
			// },
			LimiterMiddleware: limiter.SlidingWindow{}, // 使用滑动窗口算法
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(domain.Response{
					Code:    fiber.StatusTooManyRequests,
					Message: "请求过于频繁，请稍后再试",
					Data:    nil,
					Error:   nil,
				})
			},
		}))

		// 图片
		common.NewImageRoute(commonGroup.Group("image"), imageService, validate).RegisterRoutes()
		// 账号
		common.NewAccountRoute(commonGroup.Group("account"), userService, validate, privateKey).RegisterRoutes()
		// 分类
		common.NewCategoryRoute(commonGroup.Group("category"), categoryService, validate).RegisterRoutes()
		// 新闻
		common.NewArticleRoute(commonGroup.Group("article"), articleService, validate).RegisterRoutes()
		// 字典
		common.NewDictRoute(commonGroup.Group("dict"), dictService, validate).RegisterRoutes()
	}

	// 从环境变量中读取端口号，默认为 ":3000"
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // 默认端口
	}
	app.Listen(":" + port)
}
