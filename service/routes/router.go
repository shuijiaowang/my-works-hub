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

	// Serve resources files for preview (images/videos)
	r.Static("/api/resources", "./resources")

	public := r.Group("/api")
	{
		public.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
		public.POST("/example/test", apiGroup.Test)
		public.GET("/projects/all", apiGroup.ProjectApi.All)
		public.GET("/projects/:id", apiGroup.ProjectApi.Detail)
		public.GET("/projects/:id/media", apiGroup.ProjectApi.MediaList)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/login", apiGroup.AdminApi.Login)
		admin.Use(middleware.AdminAuth())
		admin.GET("/ping", apiGroup.AdminApi.Ping)
		admin.POST("/projects", apiGroup.AdminApi.CreateProject)
		admin.PUT("/projects/:id", apiGroup.AdminApi.UpdateProject)
		admin.GET("/projects/:id/media", apiGroup.AdminApi.ListProjectMedia)
		admin.POST("/projects/:id/media", apiGroup.AdminApi.UploadProjectMedia)
		admin.DELETE("/projects/:id/media/:mediaId", apiGroup.AdminApi.DeleteProjectMedia)
		admin.POST("/projects/:id/media/:mediaId/move", apiGroup.AdminApi.MoveProjectMedia)

		// zip package management (upload/update, download, delete)
		admin.POST("/projects/:id/zip", apiGroup.AdminApi.UploadProjectZip)
		admin.GET("/projects/:id/zip/:fileName", apiGroup.AdminApi.DownloadProjectZip)
		admin.DELETE("/projects/:id/zip/:fileName", apiGroup.AdminApi.DeleteProjectZip)
	}

	return r
}
