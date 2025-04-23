package admin

import (
	"cms/models/domain"
	"cms/services"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var (
	// ErrImageEmpty 图片为空
	ErrImageEmpty = errors.New("图片为空")
)

type (
	ImageRoute interface {
		RegisterRoutes()
		getImages(c *fiber.Ctx) error
		createImage(c *fiber.Ctx) error
		deleteImage(c *fiber.Ctx) error
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
	ir.app.Get("/", ir.getImages)
	ir.app.Post("/", ir.createImage)
	ir.app.Delete("/:id<int>", ir.deleteImage)
}

func (ir *imageRoute) getImages(c *fiber.Ctx) error {
	res, err := ir.imageService.GetImages()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取图片列表失败", err)
	}
	return domain.SuccessResponse(c, res, "获取图片列表成功")
}

func (ir *imageRoute) createImage(c *fiber.Ctx) error {
	formData, err := c.MultipartForm()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "获取图片失败", err)
	}

	fhs, ok := formData.File["image"]
	if !ok {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "图片为空", ErrImageEmpty)
	}

	// 判断文件夹是否存在
	absPath, err := os.Getwd()
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "获取绝对路径失败", err)
	}

	uploadPath := path.Join(absPath, "uploads")
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		// 创建文件夹
		if err := os.Mkdir(uploadPath, os.ModePerm); err != nil {
			return domain.ErrorResponse(c, fiber.StatusNotFound, "创建文件夹失败", err)
		}
	}

	res := make(domain.CreateImageResponse, 0)

	for _, fh := range fhs {
		if !strings.HasPrefix(fh.Header.Get("Content-Type"), "image/") {
			log.Infof("%s 图片格式不支持", fh.Filename)
			continue
		}

		file, err := fh.Open()
		if err != nil {
			return domain.ErrorResponse(c, fiber.StatusInternalServerError, "打开文件失败", err)
		}
		defer file.Close()

		hash := xxhash.New()
		if _, err := io.Copy(hash, file); err != nil {
			return domain.ErrorResponse(c, fiber.StatusInternalServerError, "计算哈希值失败", err)
		}

		hashSum := hash.Sum64()

		image, _ := ir.imageService.GetImageByHash(hashSum)
		if image != nil {
			log.Infof("%s 图片已存在", fh.Filename)
			res = append(res, *image)
			continue
		}

		if err = c.SaveFile(fh, path.Join(uploadPath, fmt.Sprintf("%v", hashSum))); err != nil {
			return domain.ErrorResponse(c, fiber.StatusInternalServerError, "保存图片失败", err)
		}

		params := new(domain.CreateImageParams)
		params.Title = fh.Filename
		params.Hash = hashSum

		image, err = ir.imageService.CreateImage(*params)
		if err != nil {
			return domain.ErrorResponse(c, fiber.StatusInternalServerError, "创建图片失败", err)
		}

		res = append(res, *image)
	}

	return domain.SuccessResponse(c, res, "创建图片成功")
}

func (ir *imageRoute) deleteImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return domain.ErrorResponse(c, fiber.StatusBadRequest, "参数错误", err)
	}

	if err := ir.imageService.DeleteImage(uint(id)); err != nil {
		return domain.ErrorResponse(c, fiber.StatusInternalServerError, "删除图片失败", err)
	}

	return domain.SuccessResponse(c, nil, "删除图片成功")
}
