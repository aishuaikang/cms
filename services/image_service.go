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

	// ErrImageInUse 当图片正在被使用时
	ErrImageInUse = errors.New("图片正在被使用中")
)

type (
	ImageService interface {
		GetImages() ([]*models.Image, error)
		CreateImage(image domain.CreateImageParams) (*models.Image, error)
		GetImageByHash(hash uint64) (*models.Image, error)
		GetImageById(id string) (*models.Image, error)
		DeleteImage(id string) error
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

func (s *imageService) GetImageById(id string) (*models.Image, error) {
	var image models.Image
	if err := s.db.Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (s *imageService) DeleteImage(id string) error {
	var image models.Image
	// 是否存在图片
	if err := s.db.Preload(clause.Associations).Where("id = ?", id).First(&image).Error; err != nil {
		return ErrImageNotFound
	}

	if len(image.Articles) > 0 {
		return ErrImageInUse
	}

	if err := s.db.Where("id = ?", id).Delete(&models.Image{}).Error; err != nil {
		return err
	}

	return nil
}
