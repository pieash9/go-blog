package controllers

import (
	"go-blog/database"
	"go-blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

// Create a comment on a post
func CreateComment(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, err := uuid.Parse(postIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var input CreateCommentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	comment := models.Comment{
		ID:      uuid.New(),
		Content: input.Content,
		UserID:  userID,
		PostID:  postID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully"})
}

// Get comments for a post
func GetComments(c *gin.Context) {
	postIDParam := c.Param("postId")
	var comments []models.Comment
	if err := database.DB.Preload("User").Preload("Post.User").Where("post_id = ?", postIDParam).Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// Delete comment (only by owner)
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")
	userID := c.MustGet("user_id").(uuid.UUID)

	var comment models.Comment
	if err := database.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted!"})

}
