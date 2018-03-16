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
	db, err := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	var rows *sql.Rows
	if id == "0" {
		rows, _ = db.Query("select Id,Name,'partition' as Type from partition")
	} else {
		rows, _ = db.Query("select Id,Name,'folder' as Type from folder where partitionId = ? and FolderId=0 union all select Id,Name,Type from [file] where partitionId = ? and FolderId=0", id, id)
	}

	type partition struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var pars []partition
	var par partition
	for rows.Next() {
		rows.Scan(&par.Id, &par.Name, &par.Type)
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

	db.Exec("insert into partition (name) values(?)", name)

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

	os.Remove(filepath.Join(config.ServerConfig.FilePanDir, name))

	db.Exec("delete from partition where id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
