package api

import (
	"htmlhub/config"
	"htmlhub/middleware"
	"htmlhub/util/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminApi struct{}

type AdminLoginRequest struct {
	Password string `json:"password"`
}

// Login 管理端登录：
// - 首次请求通常走 /admin/login?token=xxx（你会手动补充 token 后缀）
// - body 只做密码校验
// - 通过后返回 ok（具体会话逻辑后续再加）
func (a *AdminApi) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	expectedPwd := strings.TrimSpace(config.AppConfig.Admin.Password)
	if expectedPwd == "" {
		response.FailWithMessage("admin password 未配置", c)
		return
	}
	if strings.TrimSpace(req.Password) != expectedPwd {
		response.NoAuth("密码错误", c)
		return
	}

	// 同时要求携带固定 token（query/header 均可）
	if middleware.GetAdminToken(c) == "" {
		response.NoAuth("缺少 admin token", c)
		return
	}
	// token 正确性由 AdminAuth 统一处理，这里直接复用一次校验，避免漏配
	expectedToken := strings.TrimSpace(config.AppConfig.Admin.Token)
	if middleware.GetAdminToken(c) != expectedToken {
		response.NoAuth("admin token 校验失败", c)
		return
	}

	response.OkWithMessage("登录成功", c)
}

func (a *AdminApi) Ping(c *gin.Context) {
	response.OkWithMessage("admin pong", c)
}
