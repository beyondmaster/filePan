package book

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/beyondmaster/filePan/controller/common"
	"github.com/beyondmaster/filePan/model"
	"github.com/beyondmaster/filePan/utils"
	"github.com/gin-gonic/gin"
)

// Save 保存图书（创建或更新）
func Save(c *gin.Context, isEdit bool) {
	SendErrJSON := common.SendErrJSON
	var bookData model.Book
	if err := c.ShouldBindJSON(&bookData); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	if bookData.ContentType != model.ContentTypeMarkdown && bookData.ContentType != model.ContentTypeHTML {
		SendErrJSON("无效的图书格式", c)
		return
	}

	bookData.Name = utils.AvoidXSS(bookData.Name)
	bookData.Name = strings.TrimSpace(bookData.Name)

	bookData.Content = strings.TrimSpace(bookData.Content)
	bookData.HTMLContent = strings.TrimSpace(bookData.HTMLContent)

	if bookData.HTMLContent != "" {
		bookData.HTMLContent = utils.AvoidXSS(bookData.HTMLContent)
	}

	if bookData.Name == "" {
		SendErrJSON("图书名称不能为空", c)
		return
	}

	if utf8.RuneCountInString(bookData.Name) > model.MaxNameLen {
		msg := "图书名称不能超过" + strconv.Itoa(model.MaxNameLen) + "个字符"
		SendErrJSON(msg, c)
		return
	}

	var theContent string
	if bookData.ContentType == model.ContentTypeHTML {
		theContent = bookData.HTMLContent
	} else {
		theContent = bookData.Content
	}

	contentCount := utf8.RuneCountInString(theContent)
	if theContent == "" || contentCount <= 0 {
		SendErrJSON("图书简介不能为空", c)
		return
	}

	if contentCount > model.MaxContentLen {
		msg := "图书简介不能超过" + strconv.Itoa(model.MaxContentLen) + "个字符"
		SendErrJSON(msg, c)
		return
	}

	userInter, _ := c.Get("user")
	user := userInter.(model.User)

	var updatedBook model.Book
	if !isEdit {
		// 创建图书
		bookData.Status = model.BookUnpublish
		bookData.UserID = user.ID
		// 创建图书时，可以选择格式，之后不能修改
		if err := model.DB.Create(&bookData).Error; err != nil {
			SendErrJSON("error", c)
			return
		}
	} else {
		//更新图书
		if err := model.DB.First(&updatedBook, bookData.ID).Error; err == nil {
			updatedBook.Name = bookData.Name
			updatedBook.CoverURL = bookData.CoverURL
			updatedBook.Content = bookData.Content
			updatedBook.HTMLContent = bookData.HTMLContent
			if err := model.DB.Save(&updatedBook).Error; err != nil {
				fmt.Println(err.Error())
				SendErrJSON("error", c)
				return
			}
		} else {
			SendErrJSON("无效的图书id", c)
			return
		}
	}

	var bookJSON model.Book
	if isEdit {
		bookJSON = updatedBook
	} else {
		bookJSON = bookData
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": bookJSON,
		},
	})
}

// Create 创建图书
func Create(c *gin.Context) {
	Save(c, false)
}

// Update 更新图书
func Update(c *gin.Context) {
	Save(c, true)
}

// UpdateName 更新图书名称
func UpdateName(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var bookData model.Book
	if err := c.ShouldBindJSON(&bookData); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	bookData.Name = utils.AvoidXSS(bookData.Name)
	bookData.Name = strings.TrimSpace(bookData.Name)

	if bookData.Name == "" {
		SendErrJSON("图书名称不能为空", c)
		return
	}

	if utf8.RuneCountInString(bookData.Name) > model.MaxNameLen {
		msg := "图书名称不能超过" + strconv.Itoa(model.MaxNameLen) + "个字符"
		SendErrJSON(msg, c)
		return
	}
	var book model.Book
	if err := model.DB.First(&book, bookData.ID).Error; err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}
	book.Name = bookData.Name
	if err := model.DB.Save(&book).Error; err != nil {
		SendErrJSON("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})
}

// Publish 发布图书
func Publish(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := strconv.Atoi(c.Param("bookID"))
	if err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}
	var book model.Book
	if err := model.DB.First(&book, id).Error; err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}
	book.Status = model.BookVerifying
	if err := model.DB.Save(&book).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})
}

// List 获取图书列表
func List(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var books []model.Book
	if err := model.DB.Model(&model.Book{}).Where("status != \"book_unpublish\"").Find(&books).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
	for i := 0; i < len(books); i++ {
		if err := model.DB.Model(&books[i]).Related(&books[i].User, "users").Error; err != nil {
			fmt.Println(err.Error())
			SendErrJSON("error.", c)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"books": books,
		},
	})
}

// Info 获取图书信息
func Info(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}

	var book model.Book
	if err := model.DB.First(&book, id).Error; err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}

	if c.Query("f") != "md" {
		if book.ContentType == model.ContentTypeMarkdown {
			book.HTMLContent = utils.MarkdownToHTML(book.Content)
		} else if book.ContentType == model.ContentTypeHTML {
			book.HTMLContent = utils.AvoidXSS(book.HTMLContent)
		} else {
			book.HTMLContent = utils.MarkdownToHTML(book.Content)
		}
		book.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})
}

// Chapters 获取图书的所有章节
func Chapters(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := strconv.Atoi(c.Param("bookID"))
	if err != nil {
		SendErrJSON("错误的图书id", c)
		return
	}
	var chapters []model.BookChapter
	if err := model.DB.Model(&model.BookChapter{}).Where("book_id = ?", id).Order("created_at desc").Find(&chapters).Error; err != nil {
		fmt.Println(err)
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapters": chapters,
		},
	})
}

// Chapter 查询章节
func Chapter(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := strconv.Atoi(c.Param("chapterID"))
	if err != nil {
		SendErrJSON("错误的章节id", c)
		return
	}
	var chapter model.BookChapter
	if err := model.DB.First(&chapter, id).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}

	if c.Query("f") != "md" {
		if chapter.ContentType == model.ContentTypeMarkdown {
			chapter.HTMLContent = utils.MarkdownToHTML(chapter.Content)
		} else if chapter.ContentType == model.ContentTypeHTML {
			chapter.HTMLContent = utils.AvoidXSS(chapter.HTMLContent)
		} else {
			chapter.HTMLContent = utils.MarkdownToHTML(chapter.Content)
		}
		chapter.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}

// CreateChapter 创建图书的章节
func CreateChapter(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	type ReqData struct {
		Name     string `json:"name" binding:"required,min=1,max=100"`
		ParentID uint   `json:"parentID"`
		BookID   uint   `json:"bookID"`
	}
	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err)
		SendErrJSON("参数无效", c)
		return
	}
	reqData.Name = utils.AvoidXSS(reqData.Name)
	reqData.Name = strings.TrimSpace(reqData.Name)
	if reqData.Name == "" {
		SendErrJSON("章节名称不能为空", c)
		return
	}

	var chapter model.BookChapter
	chapter.Name = reqData.Name
	chapter.ParentID = reqData.ParentID
	chapter.BookID = reqData.BookID

	if chapter.ParentID != model.NoParent {
		var parentChapter model.BookChapter
		if err := model.DB.First(&parentChapter, chapter.ParentID).Error; err != nil {
			SendErrJSON("无效的parentID", c)
			return
		}
	}

	var book model.Book
	if err := model.DB.First(&book, chapter.BookID).Error; err != nil {
		SendErrJSON("无效的bookID", c)
		return
	}

	chapter.ContentType = book.ContentType

	if err := model.DB.Create(&chapter).Error; err != nil {
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}

// DeleteChapter 删除图书的章节
func DeleteChapter(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := strconv.Atoi(c.Param("chapterID"))
	if err != nil {
		SendErrJSON("错误的章节id", c)
		return
	}
	var sql = "DELETE FROM book_chapters WHERE id = ? OR parent_id = ?"
	if err := model.DB.Exec(sql, id, id).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": id,
		},
	})
}

// UpdateChapterContent 更新图书的章节内容
func UpdateChapterContent(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	type ReqData struct {
		ID          uint   `json:"chapterID"`
		Content     string `json:"content"`
		HTMLContent string `json:"htmlContent"`
	}
	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err)
		SendErrJSON("参数无效", c)
		return
	}

	reqData.Content = strings.TrimSpace(reqData.Content)
	reqData.HTMLContent = strings.TrimSpace(reqData.HTMLContent)

	if reqData.HTMLContent != "" {
		reqData.HTMLContent = utils.AvoidXSS(reqData.HTMLContent)
	}

	var chapter model.BookChapter
	if err := model.DB.First(&chapter, reqData.ID).Error; err != nil {
		SendErrJSON("错误的章节id", c)
		return
	}
	chapter.Content = reqData.Content
	chapter.HTMLContent = reqData.HTMLContent

	if err := model.DB.Save(&chapter).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": reqData.ID,
		},
	})
}

// UpdateChapterName 更新图书的章节的名称
func UpdateChapterName(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	type ReqData struct {
		ID   uint   `json:"id"`
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err)
		SendErrJSON("参数无效", c)
		return
	}

	reqData.Name = utils.AvoidXSS(reqData.Name)
	reqData.Name = strings.TrimSpace(reqData.Name)
	if reqData.Name == "" {
		SendErrJSON("章节名称不能为空", c)
		return
	}
	var chapter model.BookChapter
	if err := model.DB.First(&chapter, reqData.ID).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("无效的章节id", c)
		return
	}
	chapter.Name = reqData.Name
	if err := model.DB.Save(&chapter).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}
