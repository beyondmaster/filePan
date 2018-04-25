package partition

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"database/sql"
	"filePan/config"
	"filePan/model"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	id := c.Param("id")
	sort := c.Param("sort")
	db, err := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	orderBy := ""

	var rows *sql.Rows
	userId, _ := c.Get("userId")
	if id == "0" {
		if sort == "1" {
			orderBy = "IsTop desc,p.id"
		} else {
			orderBy = "Name"
		}
		rows, err = db.Query("select Id,Name,'partition' as Type,AuthLevel from partition p left join authority a on p.Id = a.partitionId where a.userid = ? order by "+orderBy, userId)
	} else {
		switch sort {
		case "1":
			orderBy = "CreatedTime"
		case "2":
			orderBy = "Name"
		case "3":
			orderBy = "Size"
		}
		rows, err = db.Query("select Id,Name,Type,AuthLevel from (select Id,Name,'folder' as Type,AuthLevel,CreatedTime,0 as Size from folder f left join authority a on f.partitionId = a.partitionId where a.userid = ? and f.partitionId = ? and FolderId=0 union all select Id,Name,Type,AuthLevel,CreatedTime,Size from [file] f left join authority a on f.partitionId = a.partitionId where a.userid = ? and f.partitionId = ? and FolderId=0) f order by "+orderBy, userId, id, userId, id)
		if err != nil {
			fmt.Println(err)
		}
	}

	type partition struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		AuthLevel int    `json:"authLevel"`
		FolderId  int    `json:"folderId"`
	}
	var pars []partition
	var par partition
	for rows.Next() {
		rows.Scan(&par.Id, &par.Name, &par.Type, &par.AuthLevel)
		pars = append(pars, par)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  pars,
	})
}

func Add(c *gin.Context) {
	name := c.Param("name")

	path := filepath.Join(config.ServerConfig.FilePanDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	//文件名存在判断
	userId, _ := c.Get("userId")
	result, _ := db.Exec("insert into partition (name,createdby) values(?,?)", name, userId)
	id, _ := result.LastInsertId()
	db.Exec("insert into authority values(?,?,3)", userId, id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")
	//文件名存在判断

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var oldName string
	row := db.QueryRow("select Name from partition where id=?", id)
	row.Scan(&oldName)

	os.Rename(filepath.Join(config.ServerConfig.FilePanDir, oldName), filepath.Join(config.ServerConfig.FilePanDir, name))

	db.Exec("update partition set name = ? where id = ?", name, id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var name string
	row := db.QueryRow("select Name from partition where id=?", id)
	row.Scan(&name)

	os.RemoveAll(filepath.Join(config.ServerConfig.FilePanDir, name))

	db.Exec("delete from partition where id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Top(c *gin.Context) {
	id := c.Param("id")
	isTop := c.Param("isTop")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var sSql string
	if isTop == "1" {
		sSql = "update partition set IsTop=0 where IsTop=1;update partition set IsTop=1 where id = ?"
	} else {
		sSql = "update partition set IsTop=0 where id = ?"
	}
	db.Exec(sSql, id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Info(c *gin.Context) {
	id := c.Param("id")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var name string
	var isTop int
	row := db.QueryRow("select Name,IsTop from partition where id=?", id)
	row.Scan(&name, &isTop)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"name":  name,
		"isTop": isTop,
	})
}
