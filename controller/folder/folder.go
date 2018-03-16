package folder

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"filePan/config"
	"filePan/model"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	id := c.Param("id")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	rows, _ := db.Query("select Id,Name,'folder' as Type from folder where folderId = ? union all select Id,Name,Type from [file] where folderId = ?", id, id)

	type file struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var files []file
	var f file
	for rows.Next() {
		rows.Scan(&f.Id, &f.Name, &f.Type)
		files = append(files, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  files,
	})
}

func Add(c *gin.Context) {
	partitionId := c.Param("partitionId")
	folderId := c.Param("folderId")
	name := c.Param("name")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var partitionName string
	var folderPath string
	var folderName string
	var path string
	var parentDir string

	if folderId == "0" {
		row := db.QueryRow("select Name from partition where id=?", partitionId)
		row.Scan(&partitionName)
		parentDir = partitionName
		path = filepath.Join(parentDir, name)
	} else {
		row := db.QueryRow("select PartitionId,Path,Name from folder where id=?", folderId)
		row.Scan(&partitionId, &folderPath, &folderName)
		parentDir = filepath.Join(folderPath, folderName)
		path = filepath.Join(parentDir, name)
	}
	path = filepath.Join(config.ServerConfig.FilePanDir, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	} else {
		//提示
	}

	db.Exec("insert into folder (Name,Path,PartitionID,FolderID) values(?,?,?,?)", name, parentDir, partitionId, folderId)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var folderPath string
	var folderName string
	row := db.QueryRow("select Path,Name from folder where id=?", id)
	row.Scan(&folderPath, &folderName)
	oldPath := filepath.Join(config.ServerConfig.FilePanDir, folderPath, folderName)
	path := filepath.Join(config.ServerConfig.FilePanDir, folderPath, name)
	//文件名存在判断

	os.Rename(oldPath, path)

	db.Exec("update folder set name = ? where id = ?", name, id)

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

	var folderPath string
	var folderName string
	row := db.QueryRow("select Path,Name from folder where id=?", id)
	row.Scan(&folderPath, &folderName)
	path := filepath.Join(config.ServerConfig.FilePanDir, folderPath, folderName)

	os.Remove(path)

	db.Exec("delete from folder where id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
