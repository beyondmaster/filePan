package authority

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"filePan/config"
	"filePan/controller/common"
	"filePan/model"
	"github.com/gin-gonic/gin"
)

type AuthInfo struct {
	PartitionId string     `json:"partitionId"`
	UserList    []UserInfo `json:"user"`
}

type UserInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	AuthLevel string `json:"authlevel"`
	Oper      string `json:"oper"`
}

func UpdateMember(c *gin.Context) {
	var authInfo AuthInfo
	if err := c.ShouldBindJSON(&authInfo); err != nil {
		common.SendErrJSON(err.Error(), c)
		return
	}
	// fmt.Println(authInfo)
	sSql := []string{}
	partitionId := authInfo.PartitionId
	for _, user := range authInfo.UserList {
		if user.Oper == "add" {
			sSql = append(sSql, fmt.Sprintf("insert into authority values('%s','%s',1)", user.Id, partitionId))
		} else {
			sSql = append(sSql, fmt.Sprintf("delete from authority where userid='%s' and partitionid='%s'", user.Id, partitionId))
		}
		sSql = append(sSql, fmt.Sprintf("if not exists (select 1 from userinfo where id = '%s') insert into userinfo values('%s','%s','%s')", user.Id, user.Id, user.Name, user.Avatar))
	}

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	fmt.Println(strings.Join(sSql, ";"))
	if _, err := db.Exec(strings.Join(sSql, ";")); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func UpdateAuth(c *gin.Context) {
	var authInfo AuthInfo
	if err := c.ShouldBindJSON(&authInfo); err != nil {
		common.SendErrJSON(err.Error(), c)
		return
	}
	// fmt.Println(authInfo)
	sSql := []string{}
	partitionId := authInfo.PartitionId
	for _, user := range authInfo.UserList {
		sSql = append(sSql, fmt.Sprintf("update authority set authlevel=%s where userid='%s' and partitionid=%s", user.AuthLevel, user.Id, partitionId))
	}

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	// fmt.Println(strings.Join(sSql, ";"))
	if _, err := db.Exec(strings.Join(sSql, ";")); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Member(c *gin.Context) {
	id := c.Param("id")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	rows, _ := db.Query("select Id,Name,Avatar,AuthLevel from authority a left join userinfo b on a.userId = b.id where partitionId = ?", id)

	type member struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Avatar    string `json:"avatar"`
		AuthLevel int    `json:"authlevel"`
	}
	var members []member
	var mem member
	for rows.Next() {
		rows.Scan(&mem.Id, &mem.Name, &mem.Avatar, &mem.AuthLevel)
		members = append(members, mem)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  members,
	})
}
