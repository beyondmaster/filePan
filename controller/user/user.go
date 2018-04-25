package user

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"filePan/config"
	"filePan/controller/common"
	"filePan/model"
	"filePan/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/chanxuehong/wechat.v1/corp/addresslist"
	"gopkg.in/chanxuehong/wechat.v1/corp/agent"
)

func Change(c *gin.Context) {
	userId := c.Param("name")
	aesEnc := utils.NewEncrypt(config.ServerConfig.TokenSecret)
	arrEncrypt, _ := aesEnc.Encrypt(userId)
	token := base64.StdEncoding.EncodeToString(arrEncrypt)
	c.SetCookie("token", token, config.ServerConfig.TokenMaxAge, "/", "", false, true)
}

func Info(c *gin.Context) {
	client := agent.NewClient(config.AccessTokenServer, nil)
	agent, err := client.AgentInfo(int64(config.WeiXinConfig.AgentId))
	if err != nil {
		common.SendErrJSON(err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  agent,
	})
}

func List(c *gin.Context) {
	departmentId, _ := strconv.ParseInt(c.Param("departmentId"), 10, 64)
	client := addresslist.NewClient(config.AccessTokenServer, nil)

	deps, err := client.DepartmentList(departmentId)
	if err != nil {
		common.SendErrJSON(err.Error(), c)
		return
	}

	if len(deps) > 0 {
		user, err := client.UserList(deps[0].Id, true, 0)
		if err != nil {
			common.SendErrJSON(err.Error(), c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"errNo":      model.ErrorCode.SUCCESS,
			"msg":        "success",
			"user":       user,
			"department": deps[1:],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"errNo":      model.ErrorCode.SUCCESS,
			"msg":        "success",
			"user":       gin.H{},
			"department": deps,
		})
	}
}
