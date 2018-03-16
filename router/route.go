package router

import (
	"filePan/config"
	"filePan/controller/file"
	"filePan/controller/folder"
	"filePan/controller/partition"
	"filePan/middleware"
	"filePan/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route 路由
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"errNo": model.ErrorCode.SUCCESS,
				"msg":   "success",
				"data":  gin.H{},
			})
		})

		api.GET("/partition/list/:id", partition.List)
		api.GET("/partition/add/:name", partition.Add)
		api.GET("/partition/update/:id/:name", partition.Update)
		api.GET("/partition/delete/:id", partition.Delete)

		api.GET("/folder/list/:id", folder.List)
		api.GET("/folder/add/:partitionId/:folderId/:name", folder.Add)
		api.GET("/folder/update/:id/:name", folder.Update)
		api.GET("/folder/delete/:id", folder.Delete)

		api.POST("/file/upload", file.Upload)
		api.GET("/file/update/:id/:name", file.Update)
		api.GET("/file/delete/:id", file.Delete)
	}

}
