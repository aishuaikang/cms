package common

import (
	"cms/models/domain"
	"cms/services"
	"fmt"
	"os"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ImageRoute interface {
		RegisterRoutes()
		downloadImageById(c *fiber.Ctx) error
	}
	imageRoute struct {
		app          fiber.Router
		imageService services.ImageService
		validator    *validator.Validate
	}
)

func NewImageRoute(app fiber.Router, imageService services.ImageService, validator *validator.Validate) ImageRoute {
	return &imageRoute{
		app:          app,
		imageService: imageService,
		validator:    validator,
	}
}

func (ir *imageRoute) RegisterRoutes() {
	ir.app.Get("/download/:id", ir.downloadImageById)

}

func (ir *imageRoute) downloadImageById(c *fiber.Ctx) error {
	id := c.Params("id")
	image, err := ir.imageService.GetImageById(id)
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusNotFound, "图片不存在", err)
	}

	// 判断文件夹是否存在
	absPath, err := os.Getwd()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取绝对路径失败", err)
	}

	uploadPath := path.Join(absPath, "uploads")
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		return domain.ErrorResponse(c, fiber.StatusNotFound, "uploads 文件夹不存在", err)
	}

	return c.Download(path.Join(uploadPath, fmt.Sprint(image.Hash)), image.Title)
}
