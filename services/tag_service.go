package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	//  标签名已存在
	ErrTagAlreadyExists = errors.New("标签名已存在")
	// 标签不存在
	ErrTagNotFound = errors.New("标签不存在")
	// 标签正在被文章使用中
	ErrTagInUseByArticle = errors.New("标签正在被文章使用中")
)

type (
	TagService interface {
		GetTags() ([]*models.Tag, error)
		CreateTag(params domain.CreateTagParams) error
		UpdateTag(id uuid.UUID, params domain.UpdateTagParams) error
		DeleteTag(id uuid.UUID) error
	}
	tagService struct {
		db *gorm.DB
	}
)

func NewTagService(db *gorm.DB) TagService {
	return &tagService{
		db: db,
	}
}
func (s *tagService) GetTags() ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := s.db.Preload(clause.Associations).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
func (s *tagService) CreateTag(params domain.CreateTagParams) error {
	// 检查标签名称是否已存在
	if err := s.db.Where("name = ?", params.Name).First(&models.Tag{}).Error; err == nil {
		return ErrTagAlreadyExists
	}

	tagModel := &models.Tag{
		Name: params.Name,
	}

	if params.Description != nil {
		tagModel.Description = *params.Description
	}

	return s.db.Create(tagModel).Error
}
func (s *tagService) UpdateTag(id uuid.UUID, params domain.UpdateTagParams) error {
	tag := new(models.Tag)

	// 检查标签是否存在
	if err := s.db.Where("id = ?", id).First(tag).Error; err != nil {
		return ErrTagNotFound
	}

	if params.Name != nil && tag.Name != *params.Name {
		// 检查标签名称是否已存在
		if err := s.db.Where("name = ?", *params.Name).First(&models.Tag{}).Error; err == nil {
			return ErrTagAlreadyExists
		}
		tag.Name = *params.Name
	}

	if params.Description != nil && tag.Description != *params.Description {
		tag.Description = *params.Description
	}

	return s.db.Save(tag).Error
}
func (s *tagService) DeleteTag(id uuid.UUID) error {
	tag := new(models.Tag)

	// 检查标签是否存在
	if err := s.db.Preload(clause.Associations).Where("id = ?", id).First(tag).Error; err != nil {
		return ErrTagNotFound
	}

	// 检查标签是否被使用
	if len(tag.Articles) > 0 {
		return ErrTagInUseByArticle
	}

	return s.db.Delete(tag).Error
}
