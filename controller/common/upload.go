package common

import (
	"errors"
	"fmt"
	// "mime"
	"database/sql"
	"net/http"
	// "os"
	// "strings"

	"filePan/config"
	"filePan/model"
	"github.com/gin-gonic/gin"
)

// Upload 文件上传
func Upload(c *gin.Context) (map[string]interface{}, error) {
	// file, err := c.FormFile("upFile")

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

	if folderId == "0" {
		row := db.QueryRow("select Name from partition where id=?", partitionId)
		row.Scan(&partitionName)
		path = partitionName
	} else {
		row := db.QueryRow("select Path,Name from folder where id=?", folderId)
		row.Scan(&folderPath, &folderName)
		path = folderPath + "/" + folderName
	}

	path = config.ServerConfig.FilePanDir + "/" + path

	for _, file := range files {

		// if err != nil {
		// 	return nil, errors.New("参数无效")
		// }

		var filename = file.Filename
		// var index = strings.LastIndex(filename, ".")

		// if index < 0 {
		// 	return nil, errors.New("无效的文件名")
		// }

		// var ext = filename[index:]
		// if len(ext) == 1 {
		// 	return nil, errors.New("无效的扩展名")
		// }
		// var mimeType = mime.TypeByExtension(ext)

		// if mimeType == "" {
		// 	return nil, errors.New("无效的图片类型")
		// }

		// imgUploadedInfo := model.GenerateImgUploadedInfo(ext)

		// fmt.Println(imgUploadedInfo.UploadDir)

		// if err := os.MkdirAll(imgUploadedInfo.UploadDir, 0777); err != nil {
		// 	fmt.Println(err.Error())
		// 	return nil, errors.New("error")
		// }

		if err := c.SaveUploadedFile(file, path+"/"+filename); err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("error1")
		}
	}

	// image := model.Image{
	// 	Title:        imgUploadedInfo.Filename,
	// 	OrignalTitle: filename,
	// 	URL:          imgUploadedInfo.ImgURL,
	// 	Width:        0,
	// 	Height:       0,
	// 	Mime:         mimeType,
	// }

	// if err := model.DB.Create(&image).Error; err != nil {
	// 	fmt.Println(err.Error())
	// 	return nil, errors.New("image error")
	// }

	return map[string]interface{}{
		"id":       "",
		"url":      "",
		"title":    "", //新文件名
		"original": "", //原始文件名
		"type":     "", //文件类型
	}, nil
}

// UploadHandler 文件上传
func UploadHandler(c *gin.Context) {
	data, err := Upload(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.ERROR,
			"msg":   err.Error(),
			"data":  gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}
