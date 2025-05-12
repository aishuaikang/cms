package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	// 字典名称已存在
	ErrDictNameAlreadyExists = errors.New("字典名称已存在")
	// 字典code已存在
	ErrDictCodeAlreadyExists = errors.New("字典code已存在")
	// 字典不存在
	ErrDictNotFound = errors.New("字典不存在")
	// 字典code不存在
	ErrDictCodeNotFound = errors.New("字典code不存在")
)

type (
	DictService interface {
		GetDicts() ([]*models.Dict, error)
		CreateDict(params domain.CreateDictParams) error
		UpdateDict(id uuid.UUID, params domain.UpdateDictParams) error
		DeleteDict(id uuid.UUID) error
		GetDictExtraByCodeWithCache(code string) (string, error)
		GetSubDictsByCodeWithCache(code string) ([]*models.Dict, error)
		GetDictByCodeWithCache(code string) (*models.Dict, error)
	}
	dictService struct {
		db *gorm.DB
	}
)

func NewDictService(db *gorm.DB) DictService {
	return &dictService{db: db}
}

func (s *dictService) GetDicts() ([]*models.Dict, error) {
	var dicts []*models.Dict
	if err := s.db.Order("created_at DESC").Find(&dicts).Error; err != nil {
		return nil, err
	}
	return dicts, nil
}
func (s *dictService) CreateDict(params domain.CreateDictParams) error {
	// 检查字典名称是否已存在
	if err := s.db.Where("name = ?", params.Name).First(&models.Dict{}).Error; err == nil {
		return ErrDictNameAlreadyExists
	}

	// 检查字典code是否已存在
	if err := s.db.Where("code = ?", params.Code).First(&models.Dict{}).Error; err == nil {
		return ErrDictCodeAlreadyExists
	}

	dictModel := &models.Dict{
		Name:        params.Name,
		Code:        params.Code,
		Extra:       params.Extra,
		Description: params.Description,
		ParentID:    params.ParentID,
	}

	if params.ImageID != nil {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", *params.ImageID).First(&models.Image{}).Error; err != nil {
			return ErrImageNotFound
		}
		dictModel.ImageID = params.ImageID
	}

	return s.db.Create(&dictModel).Error
}

func (s *dictService) UpdateDict(id uuid.UUID, params domain.UpdateDictParams) error {
	dict := new(models.Dict)
	// 检查字典是否存在
	if err := s.db.Where("id = ?", id).First(dict).Error; err != nil {
		return ErrDictNotFound
	}

	if params.Name != nil && dict.Name != *params.Name {
		// 检查字典名称是否已存在
		if err := s.db.Where("name = ?", *params.Name).First(&models.Dict{}).Error; err == nil {
			return ErrDictNameAlreadyExists
		}
		dict.Name = *params.Name
	}

	if params.Code != nil && dict.Code != *params.Code {
		// 检查字典code是否已存在
		if err := s.db.Where("code = ?", *params.Code).First(&models.Dict{}).Error; err == nil {
			return ErrDictCodeAlreadyExists
		}
		dict.Code = *params.Code
	}

	if params.Extra != nil && dict.Extra != *params.Extra {
		dict.Extra = *params.Extra
	}

	if params.ImageID != nil && dict.ImageID != params.ImageID {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", *params.ImageID).First(&models.Image{}).Error; err != nil {
			return ErrImageNotFound
		}
		dict.ImageID = params.ImageID
	}

	if params.Description != nil && dict.Description != *params.Description {
		dict.Description = *params.Description
	}

	// if params.ParentID != nil {
	// 	dictModel.ParentID = params.ParentID
	// }

	return s.db.Save(dict).Error
}

func (s *dictService) DeleteDict(id uuid.UUID) error {
	dict := new(models.Dict)
	// 检查字典是否存在
	if err := s.db.Where("id = ?", id).First(dict).Error; err != nil {
		return ErrDictNotFound
	}

	// 检查字典是否有子字典
	if err := s.db.Where("parent_id = ?", id).First(&models.Dict{}).Error; err == nil {
		return errors.New("字典有子字典，无法删除")
	}

	return s.db.Delete(dict).Error
}

func (s *dictService) GetDictExtraByCodeWithCache(code string) (string, error) {
	dict := new(models.Dict)
	if err := s.db.Where("code = ?", code).First(dict).Error; err != nil {
		return "", err
	}
	return dict.Extra, nil
}

func (s *dictService) GetDictByCodeWithCache(code string) (*models.Dict, error) {
	dict := new(models.Dict)
	if err := s.db.Where("code = ?", code).First(dict).Error; err != nil {
		return nil, err
	}
	return dict, nil
}

func (s *dictService) GetSubDictsByCodeWithCache(code string) ([]*models.Dict, error) {
	dict := new(models.Dict)
	// 检查字典是否存在
	if err := s.db.Where("code = ?", code).First(dict).Error; err != nil {
		return nil, ErrDictCodeNotFound
	}

	// 获取字典的子字典列表
	var subDicts []*models.Dict
	if err := s.db.Where("parent_id = ?", dict.ID).Find(&subDicts).Error; err != nil {
		return nil, err
	}
	if len(subDicts) == 0 {
		return nil, ErrDictNotFound
	}
	// 将子字典列表添加到字典列表中
	return subDicts, nil
}
