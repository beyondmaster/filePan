package router

import (
	"filePan/config"
	"filePan/controller/authority"
	"filePan/controller/file"
	"filePan/controller/folder"
	"filePan/controller/partition"
	"filePan/controller/user"
	"filePan/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route 路由
func Route(router *gin.Engine) {

	router.Static("/assets", "website/assets")
	router.LoadHTMLGlob("website/*.html")

	router.Use(middleware.RefreshTokenCookie)
	router.Use(middleware.AuthRequired)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/detail.html", func(c *gin.Context) {
		authLevel, _ := c.Get("authLevel")
		c.HTML(http.StatusOK, "detail.html", map[string]interface{}{
			"authLevel": authLevel,
		})
	})
	router.GET("/file-detail.html", func(c *gin.Context) {
		authLevel, _ := c.Get("authLevel")
		c.HTML(http.StatusOK, "file-detail.html", map[string]interface{}{
			"authLevel": authLevel,
		})
	})
	router.GET("/accredit.html", middleware.AdminRequired, func(c *gin.Context) {
		c.HTML(http.StatusOK, "accredit.html", nil)
	})
	router.GET("/set-accredit.html", middleware.AdminRequired, func(c *gin.Context) {
		c.HTML(http.StatusOK, "set-accredit.html", nil)
	})
	router.GET("/choose-member.html", middleware.AdminRequired, func(c *gin.Context) {
		c.HTML(http.StatusOK, "choose-member.html", nil)
	})

	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix)
	{
		api.GET("/partition/list/:id/:sort", partition.List)
		api.GET("/partition/add/:name", partition.Add)
		api.GET("/partition/update/:id/:name", middleware.AdminRequired, partition.Update)
		api.GET("/partition/delete/:id", middleware.AdminRequired, partition.Delete)
		api.GET("/partition/top/:id/:isTop", middleware.AdminRequired, partition.Top)
		api.GET("/partition/info/:id", middleware.AdminRequired, partition.Info)

		api.GET("/folder/list/:id/:sort", folder.List)
		api.GET("/folder/add/:partitionId/:folderId/:name", middleware.EditorRequired, folder.Add)
		api.GET("/folder/update/:id/:name", middleware.EditorRequired, folder.Update)
		api.GET("/folder/delete/:id", middleware.EditorRequired, folder.Delete)
		api.GET("/folder/tree/:partitionId/:folderId", folder.Tree)
		api.GET("/folder/move/:id/:folderId", middleware.EditorRequired, folder.Move)

		api.POST("/file/upload", middleware.EditorRequired, file.Upload)
		api.GET("/file/update/:id/:name", middleware.EditorRequired, file.Update)
		api.GET("/file/delete/:id", middleware.EditorRequired, file.Delete)
		api.GET("/file/info/:id", file.Info)
		api.GET("/file/move/:id/:folderId", middleware.EditorRequired, file.Move)

		api.GET("/user/change/:name", user.Change)
		api.GET("/user/list/:departmentId", middleware.AdminRequired, user.List)

		api.GET("/authority/member/:id", middleware.AdminRequired, authority.Member)
		api.POST("/authority/updateMember", middleware.AdminRequired, authority.UpdateMember)
		api.POST("/authority/updateAuth", middleware.AdminRequired, authority.UpdateAuth)
	}
}
