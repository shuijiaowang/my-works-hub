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
		public.GET("/projects/:folderName", apiGroup.ProjectApi.Detail)
		public.GET("/projects/:folderName/media", apiGroup.ProjectApi.MediaList)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/login", apiGroup.AdminApi.Login)
		admin.Use(middleware.AdminAuth())
		admin.GET("/ping", apiGroup.AdminApi.Ping)
		admin.POST("/projects", apiGroup.AdminApi.CreateProject)
		admin.PUT("/projects/:folderName", apiGroup.AdminApi.UpdateProject)
		admin.GET("/projects/:folderName/media", apiGroup.AdminApi.ListProjectMedia)
		admin.POST("/projects/:folderName/media", apiGroup.AdminApi.UploadProjectMedia)
		admin.DELETE("/projects/:folderName/media/:mediaId", apiGroup.AdminApi.DeleteProjectMedia)
		admin.POST("/projects/:folderName/media/:mediaId/move", apiGroup.AdminApi.MoveProjectMedia)

		// zip package management (upload/update, download, delete)
		admin.POST("/projects/:folderName/zip", apiGroup.AdminApi.UploadProjectZip)
		admin.GET("/projects/:folderName/zip/:fileName", apiGroup.AdminApi.DownloadProjectZip)
		admin.DELETE("/projects/:folderName/zip/:fileName", apiGroup.AdminApi.DeleteProjectZip)
	}

	return r
}
