package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"go-chat-app/config"
	"go-chat-app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins — adjust for production
	},
}

// Connected clients
var clients = make(map[string]*websocket.Conn)

// Rate limiter structure
type RateLimiter struct {
	clients map[string][]time.Time
	mu      sync.Mutex
}

// Create new rate limiter instance
var limiter = RateLimiter{
	clients: make(map[string][]time.Time),
}

// Message structure
type WSMessage struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

// WebSocket connection
func HandleWebSocket(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse JWT
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("supersecretkey"), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract username from JWT token
	username := claims["username"].(string)

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅", username, "connected via WebSocket")
	clients[username] = conn
	defer delete(clients, username)

	// Read messages
	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			break
		}

		// 🛑 **Apply Rate Limiting**
		limiter.mu.Lock()
		timestamps := limiter.clients[username]
		now := time.Now()

		// Keep only recent timestamps within the 10-second window
		validTimestamps := []time.Time{}
		for _, t := range timestamps {
			if now.Sub(t) < 10*time.Second {
				validTimestamps = append(validTimestamps, t)
			}
		}

		// Check if user exceeded message limit
		if len(validTimestamps) >= 5 {
			conn.WriteJSON(gin.H{"error": "Rate limit exceeded. Max 5 messages per 10 seconds."})
			limiter.mu.Unlock()
			continue
		}

		// Add timestamp
		validTimestamps = append(validTimestamps, now)
		limiter.clients[username] = validTimestamps
		limiter.mu.Unlock()

		fmt.Println("Received message:", msg)

		// Save message in database
		message := models.Message{
			Sender:    username,
			Recipient: msg.Recipient,
			Content:   msg.Content,
			Timestamp: time.Now(),
		}

		// Debug log before saving
		fmt.Println("Saving message to DB:", message)

		result := config.DB.Create(&message)
		if result.Error != nil {
			fmt.Println("❌ Failed to save message:", result.Error)
		} else {
			fmt.Println("✅ Message saved to DB!")
		}

		// Deliver to recipient if online
		if recipientConn, ok := clients[msg.Recipient]; ok {
			recipientConn.WriteJSON(message)
		}
	}
}
