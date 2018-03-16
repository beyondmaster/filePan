package file

import (
	"database/sql"
	// "errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"filePan/config"
	"filePan/model"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upFile"]

	partitionId := c.PostForm("partitionId")
	folderId := c.PostForm("folderId")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var partitionName string
	var folderPath string
	var folderName string
	var path string
	var absPath string

	if folderId == "0" {
		row := db.QueryRow("select Name from partition where id=?", partitionId)
		row.Scan(&partitionName)
		path = partitionName
	} else {
		row := db.QueryRow("select PartitionId,Path,Name from folder where id=?", folderId)
		row.Scan(&partitionId, &folderPath, &folderName)
		path = filepath.Join(folderPath, folderName)
	}

	absPath = filepath.Join(config.ServerConfig.FilePanDir, path)
	sSql := []string{}
	for _, file := range files {
		if err := c.SaveUploadedFile(file, filepath.Join(absPath, file.Filename)); err != nil {
			fmt.Println(err.Error())
			// return nil, errors.New("error")
		} else {
			sSql = append(sSql, fmt.Sprintf("insert into [file] (Name,Type,Size,Path,PartitionId,FolderId) values('%s','%s',%d,'%s',%s,%s)", file.Filename, file.Filename[(strings.LastIndex(file.Filename, ".")+1):], file.Size, path, partitionId, folderId))
		}
	}
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

func Update(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var folderPath string
	var fileName string
	row := db.QueryRow("select Path,Name from [file] where id=?", id)
	row.Scan(&folderPath, &fileName)
	oldPath := filepath.Join(config.ServerConfig.FilePanDir, folderPath, fileName)
	path := filepath.Join(config.ServerConfig.FilePanDir, folderPath, name)
	//文件名存在判断

	os.Rename(oldPath, path)

	db.Exec("update [file] set name = ?,type = ? where id = ?", name, name[(strings.LastIndex(name, ".")+1):], id)

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
	var fileName string
	row := db.QueryRow("select Path,Name from [file] where id=?", id)
	row.Scan(&folderPath, &fileName)
	path := filepath.Join(config.ServerConfig.FilePanDir, folderPath, fileName)

	os.Remove(path)

	db.Exec("delete from [file] where id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
