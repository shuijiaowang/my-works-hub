package middleware

import (
	"htmlhub/config"
	"htmlhub/util/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AdminTokenQueryKey  = "token"
	AdminTokenHeaderKey = "X-Admin-Token"
)

func GetAdminToken(c *gin.Context) string {
	if t := strings.TrimSpace(c.Query(AdminTokenQueryKey)); t != "" {
		return t
	}
	if t := strings.TrimSpace(c.GetHeader(AdminTokenHeaderKey)); t != "" {
		return t
	}
	return ""
}

// AdminAuth 只保护管理端路由：要求请求携带固定 token。
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := strings.TrimSpace(config.AppConfig.Admin.Token)
		if expected == "" {
			response.NoAuth("admin token 未配置", c)
			c.Abort()
			return
		}

		got := GetAdminToken(c)
		if got == "" || got != expected {
			response.NoAuth("无权限：admin token 校验失败", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
