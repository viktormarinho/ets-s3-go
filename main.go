package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/viktormarinho/ets-s3-go/models"
	"github.com/viktormarinho/ets-s3-go/useful"
	"gorm.io/gorm"
)

type User struct {
	ID        *string `json:"id"`
	Email     *string `json:"email"`
	Sector    *string `json:"sector"`
	Edv       *string `json:"edv"`
	AvatarUrl *string `json:"avatarUrl"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	RoleName  *string `json:"roleName"`
}

type AuthResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	router.Use(AuthMiddleware)

	router.POST("/upload/file", UploadFile)

	router.GET("/retrieve/file/:id", RetrieveFile)

	router.GET("/download/file/:id", DownloadFile)

	router.Run()
}

func DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")

	var dbFile models.File

	result := models.DB.First(&dbFile, id)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(404, gin.H{
			"error": "file with given id not found",
		})
		return
	}

	filename := useful.GetFilename(dbFile.Path)

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Type", dbFile.Type)
	ctx.Header("Content-Encoding", "gzip")
	ctx.FileAttachment(dbFile.Path, filename)
}

func RetrieveFile(ctx *gin.Context) {
	id := ctx.Param("id")

	var file models.File

	result := models.DB.First(&file, id)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(404, gin.H{
			"error": "file with given id not found",
		})
		return
	}

	downloadUrl := fmt.Sprintf("http://%s/download/file/%d", ctx.Request.Host, file.ID)

	ctx.JSON(200, gin.H{
		"file":        file,
		"downloadUrl": downloadUrl,
	})
}

func UploadFile(c *gin.Context) {
	user, ok := c.MustGet("currentUser").(User)

	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "could not type cast internal context user object"})
		return
	}

	incomingFile, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "no file sent"})
		return
	}

	meta := c.Request.FormValue("meta")

	path := fmt.Sprintf("./uploads/%d-%s", time.Now().Unix(), incomingFile.Filename)

	if err := useful.CompressAndSave(incomingFile, path); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"message": "could not save uploaded file to system",
			"error":   err.Error(),
		})
		return
	}

	dbFile := models.File{Type: incomingFile.Header.Get("Content-Type"), Path: path, UserId: *user.ID, Meta: &meta}

	models.DB.Create(&dbFile)

	c.JSON(201, gin.H{
		"file": dbFile,
	})
}
