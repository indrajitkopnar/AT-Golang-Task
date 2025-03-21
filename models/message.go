package models

import (
	"time"
)

// Message struct defines the chat message model
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Sender    string    `gorm:"not null" json:"sender"`
	Recipient string    `gorm:"not null" json:"recipient"`
	Content   string    `gorm:"not null" json:"content"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}
