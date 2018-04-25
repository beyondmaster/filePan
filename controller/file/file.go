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

	userId, _ := c.Get("userId")
	absPath = filepath.Join(config.ServerConfig.FilePanDir, path)
	sSql := []string{}
	for _, file := range files {
		if err := c.SaveUploadedFile(file, filepath.Join(absPath, file.Filename)); err != nil {
			fmt.Println(err.Error())
		} else {
			sSql = append(sSql, fmt.Sprintf("insert into [file] (Name,Type,Size,Path,PartitionId,FolderId,createdby) values('%s','%s',%d,'%s',%s,%s,'%s')", file.Filename, file.Filename[(strings.LastIndex(file.Filename, ".")+1):], file.Size, path, partitionId, folderId, userId))
		}
	}
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

func Info(c *gin.Context) {
	id := c.Param("id")
	db, err := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	row := db.QueryRow("select Type,Path,Name,Size from [file] where id=?", id)

	type fileinfo struct {
		Type string `json:"type"`
		path string `json:"path"`
		Name string `json:"name"`
		Size int    `json:"size"`
		Link string `json:"link"`
	}
	var info fileinfo
	row.Scan(&info.Type, &info.path, &info.Name, &info.Size)

	info.Link = strings.Join([]string{fmt.Sprintf("http://%s", config.ServerConfig.ImgHost), info.path, info.Name}, "/")

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  info,
	})
}

func Move(c *gin.Context) {
	id := c.Param("id")
	folderId := c.Param("folderId")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var partitionName string
	var folderPath string
	var folderName string
	var path string
	var sourceFolder string
	var name string

	row := db.QueryRow("select Path,Name from [file] where id=?", id)
	row.Scan(&sourceFolder, &name)

	fmt.Println(sourceFolder, name)

	if folderId == "0" {
		row := db.QueryRow("select p.Name from [file] f left join partition p on f.partitionId = p.id where f.id = ?", id)
		row.Scan(&partitionName)
		path = partitionName
	} else {
		row := db.QueryRow("select Path,Name from folder where id=?", folderId)
		row.Scan(&folderPath, &folderName)
		path = filepath.Join(folderPath, folderName)
	}

	os.Rename(filepath.Join(config.ServerConfig.FilePanDir, sourceFolder, name), filepath.Join(config.ServerConfig.FilePanDir, path, name))
	fmt.Println(filepath.Join(config.ServerConfig.FilePanDir, sourceFolder, name), filepath.Join(config.ServerConfig.FilePanDir, path, name))

	_, err := db.Exec("update [file] set path = ?, folderId = ? where id = ?", path, folderId, id)

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
