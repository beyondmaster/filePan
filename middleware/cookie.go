package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"filePan/config"
	"filePan/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/chanxuehong/wechat.v1/corp/oauth2"
)

// RefreshTokenCookie 刷新过期时间
func RefreshTokenCookie(c *gin.Context) {
	aesEnc := utils.NewEncrypt(config.ServerConfig.TokenSecret)
	var userId string
	token, err := c.Cookie("token")
	if err != nil {
		if strings.Contains(c.Request.UserAgent(), "MicroMessenger") {
			code := c.Query("code")
			if code == "" {
				fmt.Println("redirect")
				c.Redirect(http.StatusFound, oauth2.AuthCodeURL(config.WeiXinConfig.CorpID, strings.Join([]string{"http://", config.ServerConfig.Host, c.Request.RequestURI}, ""), "snsapi_base", "state"))
				return
			} else {
				client := oauth2.NewClient(config.AccessTokenServer, nil)
				user, _ := client.UserInfo(int64(config.WeiXinConfig.AgentId), code)
				fmt.Println("get info")
				fmt.Println(user)
				userId = user.UserId
				arrEncrypt, _ := aesEnc.Encrypt(userId)
				token = base64.StdEncoding.EncodeToString(arrEncrypt)
				c.SetCookie("token", token, config.ServerConfig.TokenMaxAge, "/", "", false, true)
			}
		} else {
			userId = "WangChenXu"
			arrEncrypt, _ := aesEnc.Encrypt(userId)
			token = base64.StdEncoding.EncodeToString(arrEncrypt)
			c.SetCookie("token", token, config.ServerConfig.TokenMaxAge, "/", "", false, true)
		}
	} else {
		arrEncrypt, _ := base64.StdEncoding.DecodeString(token)
		userId, _ = aesEnc.Decrypt(arrEncrypt)
	}
	fmt.Println("userId:", userId)
	c.Set("userId", userId)
	c.Next()
}
