package domain

import "cms/models"

type (
	// 用户登录参数
	LoginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"` // 密码
		// Captcha   string `json:"captcha" validate:"required"` // 验证码
	}
	// 用户登录响应
	LoginResponse struct {
		models.User
		Token string `json:"token"` // token
	}

	// 添加用户参数
	CreateUserParams struct {
		Nickname string `json:"nickname" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		ImageID  *uint  `json:"image_id,string"`
	}

	// 更新用户参数
	UpdateUserParams struct {
		Nickname *string `json:"nickname"`
		Phone    *string `json:"phone"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		ImageID  *uint   `json:"image_id,string"`
	}
)
