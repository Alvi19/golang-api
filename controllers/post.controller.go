package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alvi19/golang-gorm-postgres/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostController {
	return PostController{DB}
}

// [...] Create Post Handler
// func (pc *PostController) CreatePost(ctx *gin.Context) {
// 	currentUser := ctx.MustGet("currentUser").(models.User)
// 	var payload *models.CreatePostRequest

// 	if err := ctx.ShouldBindJSON(&payload); err != nil {
// 		ctx.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	now := time.Now()
// 	newPost := models.Post{
// 		Title:     payload.Title,
// 		Content:   payload.Content,
// 		Image:     payload.Image,
// 		User:      currentUser.ID,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	result := pc.DB.Create(&newPost)
// 	if result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "duplicate key") {
// 			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
// 			return
// 		}
// 		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
// }

// func (pc *PostController) CreatePost(ctx *gin.Context) {
// 	currentUser := ctx.MustGet("currentUser").(models.User)

// 	// Parse the form data
// 	title := ctx.PostForm("title")
// 	content := ctx.PostForm("content")

// 	// Handle file upload
// 	file, err := ctx.FormFile("image")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Image upload failed"})
// 		return
// 	}

// 	// Save the file to a desired location
// 	imagePath := "uploads/" + file.Filename
// 	if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save image"})
// 		return
// 	}

// 	now := time.Now()
// 	newPost := models.Post{
// 		Title:     title,
// 		Content:   content,
// 		Image:     imagePath,
// 		User:      currentUser.ID,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	result := pc.DB.Create(&newPost)
// 	if result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "duplicate key") {
// 			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
// 			return
// 		}
// 		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
// }

// [Create Post]
func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Ambil data dari form
	name := ctx.PostForm("name")
	lokasi := ctx.PostForm("lokasi")

	// Periksa apakah title dan content tidak kosong
	if name == "" || lokasi == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Data harus diisi"})
		return
	}

	// Tangani upload file
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Upload gambar gagal"})
		return
	}

	// Simpan file ke lokasi yang diinginkan
	imagePath := "uploads/" + file.Filename
	if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan gambar"})
		return
	}

	now := time.Now()
	newPost := models.Post{
		Name:      name,
		Lokasi:    lokasi,
		Image:     imagePath,
		Status:    true,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data tersebut sudah ada"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}

// [...] Update Post Handler
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postID := ctx.Param("postId")
	var post models.Post

	// Cari post berdasarkan ID
	if err := pc.DB.First(&post, "id = ?", postID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Post tidak ditemukan"})
		return
	}

	// Ambil data dari form
	name := ctx.PostForm("name")
	lokasi := ctx.PostForm("lokasi")
	status := ctx.PostForm("status")

	// Periksa apakah name dan lokasi tidak kosong
	if name == "" && lokasi == "" && ctx.Request.MultipartForm == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Data harus diisi"})
		return
	}

	// Tangani upload file jika ada
	file, err := ctx.FormFile("image")
	if err == nil {
		// Simpan file ke lokasi yang diinginkan
		imagePath := "uploads/" + file.Filename
		if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan gambar"})
			return
		}
		post.Image = imagePath
	}

	// Perbarui fields yang disediakan
	if name != "" {
		post.Name = name
	}
	if lokasi != "" {
		post.Lokasi = lokasi
	}
	if status != "" {
		statusValue, err := strconv.ParseBool(status)
		if err == nil {
			post.Status = statusValue
		}
	}
	post.UpdatedAt = time.Now()

	if err := pc.DB.Save(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// [...] Get Single Post Handler
func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// [...] Get All Posts Handler
func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var posts []models.Post
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&posts)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	result := pc.DB.Delete(&models.Post{}, "id = ?", postId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that ID exists"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data deleted successfully"})
}
