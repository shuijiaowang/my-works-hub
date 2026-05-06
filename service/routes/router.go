package routes

import (
	"workhub/api"
	"workhub/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 跨域中间件（放在最前面）
	r.Use(middleware.Cors())

	apiGroup := api.ApiGroup{}

	public := r.Group("/api")
	{
		public.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
		public.POST("/example/test", apiGroup.Test)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/login", apiGroup.AdminApi.Login)
		admin.Use(middleware.AdminAuth())
		admin.GET("/ping", apiGroup.AdminApi.Ping)
	}

	return r
}
