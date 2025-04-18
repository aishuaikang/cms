package domain

type (
	// 用户登录参数
	LoginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"` // 密码
		// Captcha   string `json:"captcha" validate:"required"` // 验证码
	}
)
