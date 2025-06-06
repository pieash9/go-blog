package controllers

import (
	"go-blog/database"
	"go-blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// create new blog post
func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("user_id").(uuid.UUID)

	post := models.Post{
		ID:      uuid.New(),
		Title:   input.Title,
		Content: input.Content,
		UserID:  userId,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

// get all posts
func GetPosts(c *gin.Context) {
	var posts []models.Post
	database.DB.Preload("User").Order("created_at desc").Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// Get a single post by ID
func GetPostById(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Preload("User").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// UpdatePost updates a post by ID, only if the user is the owner
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uuid.UUID)

	var post models.Post
	if err := database.DB.First(&post, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if current user is the owner
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post"})
		return
	}

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = input.Title
	post.Content = input.Content

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// Delete a post (only by owner)
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uuid.UUID)

	var post models.Post
	if err := database.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	database.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
