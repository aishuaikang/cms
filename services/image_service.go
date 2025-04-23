package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrImageNotFound 当图片不存在时
	ErrImageNotFound = errors.New("图片不存在")

	// ErrImageInUseByArticle 当图片正在被文章使用中
	ErrImageInUseByArticle = errors.New("图片正在被文章使用中")

	// ErrImageInUseByUser 当前图片正在被用户使用中
	ErrImageInUseByUser = errors.New("图片正在被用户使用中")
)

type (
	ImageService interface {
		GetImages() ([]*models.Image, error)
		CreateImage(image domain.CreateImageParams) (*models.Image, error)
		GetImageByHash(hash uint64) (*models.Image, error)
		GetImageById(id uint) (*models.Image, error)
		DeleteImage(id uint) error
	}
	imageService struct {
		db *gorm.DB
	}
)

func NewImageService(db *gorm.DB) ImageService {
	return &imageService{db: db}
}

func (s *imageService) GetImages() ([]*models.Image, error) {
	var images []*models.Image
	if err := s.db.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (s *imageService) CreateImage(image domain.CreateImageParams) (*models.Image, error) {
	imageModel := &models.Image{
		Title: image.Title,
		Hash:  image.Hash,
	}

	if err := s.db.Create(imageModel).Error; err != nil {
		return nil, err
	}

	return imageModel, nil
}

func (s *imageService) GetImageByHash(hash uint64) (*models.Image, error) {
	var image models.Image
	if err := s.db.Where("hash = ?", hash).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (s *imageService) GetImageById(id uint) (*models.Image, error) {
	var image models.Image
	if err := s.db.Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (s *imageService) DeleteImage(id uint) error {
	image := new(models.Image)

	// 检查图片是否存在
	if err := s.db.Preload(clause.Associations).Where("id = ?", id).First(image).Error; err != nil {
		return ErrImageNotFound
	}

	// 检查图片是否正在被使用
	if len(image.Articles) > 0 {
		return ErrImageInUseByArticle
	}

	// 检查图片是否正在被用户使用
	if len(image.Users) > 0 {
		return ErrImageInUseByUser
	}

	return s.db.Delete(image).Error
}
