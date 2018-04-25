package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"filePan/config"
	"github.com/gin-gonic/gin"
)

func getAuthLevel(c *gin.Context) int {
	path := c.Request.URL.Path
	var authLevel int
	userId, _ := c.Get("userId")

	var id = c.Param("id")
	var sSql = "select AuthLevel from authority where userId=? and partitionId=?"
	switch {
	case path == "/" || path == "/index.html" || path == "/api/partition/list/0":
		return 1
	case path == "/accredit.html" || path == "/choose-member.html" || path == "/set-accredit.html":
		id = c.Query("id")
	case path == "/detail.html":
		if c.Query("pid") != "" {
			id = c.Query("pid")
		} else if c.Query("fid") != "" {
			id = c.Query("fid")
			sSql = "select AuthLevel from authority a left join folder f on a.PartitionId = f.PartitionId where userId=? and f.Id=?"
		}
	case path == "/file-detail.html":
		id = c.Query("id")
		sSql = "select AuthLevel from authority a left join [file] f on a.PartitionId = f.PartitionId where userId=? and f.Id=?"
	case strings.HasPrefix(path, "/api/partition/list"):
		return 1
	case strings.HasPrefix(path, "/api/partition/add"):
		return 1
	case strings.HasPrefix(path, "/api/user/list"):
		return 3
	case strings.HasPrefix(path, "/api/user/change"):
		return 3
	case strings.HasPrefix(path, "/api/folder/tree"):
		return 1
	case strings.HasPrefix(path, "/api/folder/add"):
		id = c.Param("partitionId")
	case strings.HasPrefix(path, "/api/folder"):
		sSql = "select AuthLevel from authority a left join folder f on a.PartitionId = f.PartitionId where userId=? and f.Id=?"
	case path == "/api/file/upload":
		if c.PostForm("partitionId") != "" {
			id = c.PostForm("partitionId")
		} else if c.PostForm("folderId") != "" {
			id = c.PostForm("folderId")
			sSql = "select AuthLevel from authority a left join folder f on a.PartitionId = f.PartitionId where userId=? and f.Id=?"
		}
	case strings.HasPrefix(path, "/api/file"):
		sSql = "select AuthLevel from authority a left join [file] f on a.PartitionId = f.PartitionId where userId=? and f.Id=?"
	case strings.HasPrefix(path, "/api/authority/updateMember"):
		return 3
	case strings.HasPrefix(path, "/api/authority/updateAuth"):
		return 3
	}

	fmt.Println(userId, id)
	if userId != "" && id != "" {
		db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
		defer db.Close()

		row := db.QueryRow(sSql, userId, id)
		row.Scan(&authLevel)
	}

	return authLevel
}

// AuthRequired 必须是授权用户
func AuthRequired(c *gin.Context) {
	authLevel := getAuthLevel(c)
	if authLevel < 1 {
		c.String(http.StatusOK, "未授权")
		c.Abort()
		return
	}
	c.Set("authLevel", authLevel)
	c.Next()
}

// EditorRequired 必须是可编辑
func EditorRequired(c *gin.Context) {
	authLevel, _ := c.Get("authLevel")
	if authLevel.(int) < 2 {
		c.String(http.StatusOK, "未授权")
		c.Abort()
		return
	}
	c.Next()
}

// AdminRequired 必须是管理员
func AdminRequired(c *gin.Context) {
	authLevel, _ := c.Get("authLevel")
	if authLevel.(int) < 3 {
		c.String(http.StatusOK, "未授权")
		c.Abort()
		return
	}
}
