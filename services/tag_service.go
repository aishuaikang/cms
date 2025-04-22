package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	//  标签名已存在
	ErrTagAlreadyExists = errors.New("标签名已存在")
	// 标签不存在
	ErrTagNotFound = errors.New("标签不存在")
	// 当前标签正在被文章使用中
	ErrTagInUseByArticle = errors.New("当前标签正在被文章使用中")
)

type (
	TagService interface {
		GetTags() ([]*models.Tag, error)
		CreateTag(params domain.CreateTagParams) error
		UpdateTag(id string, params domain.UpdateTagParams) error
		DeleteTag(id string) error
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
	if err := s.db.Find(&tags).Error; err != nil {
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
func (s *tagService) UpdateTag(id string, params domain.UpdateTagParams) error {
	// 检查标签是否存在
	if err := s.db.Where("id = ?", id).First(&models.Tag{}).Error; err != nil {
		return ErrTagNotFound
	}

	tagModel := &models.Tag{}

	if params.Name != nil {
		// 检查标签名称是否已存在
		if err := s.db.Where("name = ?", *params.Name).First(&models.Tag{}).Error; err == nil {
			return ErrTagAlreadyExists
		}
		tagModel.Name = *params.Name
	}

	if params.Description != nil {
		tagModel.Description = *params.Description
	}

	return s.db.Model(&models.Tag{}).Where("id = ?", id).Updates(tagModel).Error
}
func (s *tagService) DeleteTag(id string) error {
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
