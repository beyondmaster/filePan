package folder

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"filePan/config"
	"filePan/model"
	"filePan/utils"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	id := c.Param("id")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	rows, _ := db.Query("select Id,Name,'folder' as Type,FolderId from folder where folderId = ? union all select Id,Name,Type,FolderId from [file] where folderId = ?", id, id)

	type file struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		FolderId int    `json:"folderId"`
	}
	var files []file
	var f file
	for rows.Next() {
		rows.Scan(&f.Id, &f.Name, &f.Type, &f.FolderId)
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
	userId, _ := c.Get("userId")
	db.Exec("insert into folder (Name,Path,PartitionID,FolderID,createdby) values(?,?,?,?,?)", name, parentDir, partitionId, folderId, userId)

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

	os.RemoveAll(path)

	db.Exec("delete from folder where id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func Tree(c *gin.Context) {
	partitionId := c.Param("partitionId")
	folderId := c.Param("folderId")

	db, _ := sql.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	defer db.Close()

	var rows *sql.Rows
	var err error
	if folderId == "-1" {
		rows, _ = db.Query("select 0 as Id,Name,case when cnt>0 then 'folder' else 'item' end as Type,Id as PartitionId from partition p left join (select PartitionId,count(*) as cnt from folder where FolderId = 0 group by PartitionId) f on p.Id = f.PartitionId where p.Id = ?", partitionId)
	} else if partitionId == "0" {
		rows, err = db.Query("select 0 as Id,p.Name,case when cnt>0 then 'folder' else 'item' end as Type,PartitionId from folder f left join partition p on f.partitionId = p.id left join (select PartitionId as Id,count(*) as cnt from folder where FolderId = 0 group by PartitionId) n on n.Id = p.Id where f.Id = ?", folderId)
	} else {
		rows, err = db.Query("select f.Id,Name,case when cnt>0 then 'folder' else 'item' end as Type,PartitionId from folder f left join (select FolderId as Id,count(*) as cnt from folder group by FolderId) n on n.Id = f.id where partitionId = ? and folderId = ?", partitionId, folderId)
	}

	if err != nil {
		fmt.Println(err)
	}

	type folder struct {
		Id          int    `json:"id"`
		Name        string `json:"title"`
		Type        string `json:"type"`
		PartitionId string `json:"partitionId"`
	}
	var folders []folder
	var f folder
	for rows.Next() {
		rows.Scan(&f.Id, &f.Name, &f.Type, &f.PartitionId)
		folders = append(folders, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  folders,
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

	row := db.QueryRow("select Path,Name from folder where id=?", id)
	row.Scan(&sourceFolder, &name)

	if folderId == "0" {
		row := db.QueryRow("select p.Name from folder f left join partition p on f.partitionId = p.id where f.id = ?", id)
		row.Scan(&partitionName)
		path = partitionName
	} else {
		row := db.QueryRow("select Path,Name from folder where id=?", folderId)
		row.Scan(&folderPath, &folderName)
		path = filepath.Join(folderPath, folderName)
	}

	utils.CopyDir(filepath.Join(config.ServerConfig.FilePanDir, sourceFolder, name), filepath.Join(config.ServerConfig.FilePanDir, path, name))
	os.RemoveAll(filepath.Join(config.ServerConfig.FilePanDir, sourceFolder, name))

	fmt.Println(filepath.Join(config.ServerConfig.FilePanDir, sourceFolder, name), filepath.Join(config.ServerConfig.FilePanDir, path, name))

	_, err := db.Exec("update folder set path = ?, folderId = ? where id = ?", path, folderId, id)

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
