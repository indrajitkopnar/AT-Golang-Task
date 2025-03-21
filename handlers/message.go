package handlers

import (
	"net/http"
	"strconv"

	"go-chat-app/config"
	"go-chat-app/models"

	"github.com/gin-gonic/gin"
)

// FetchChatHistory retrieves chat messages between two users
func FetchChatHistory(c *gin.Context) {
	user1 := c.Query("user1")
	user2 := c.Query("user2")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	var messages []models.Message

	// Query messages between two users (both directions)
	config.DB.Where(
		"(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)",
		user1, user2, user2, user1,
	).Order("timestamp DESC").Limit(limit).Offset(offset).Find(&messages)

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
