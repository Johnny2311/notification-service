package entity

import (
	"notification-service/enum"
	"time"
)

type Notification struct {
	ID      int64
	Date    time.Time
	Name    string
	Method  enum.DeliveryMethod
	Type    enum.NotificationType
	Email   string
	UserID  string
	Content string
	Seen    bool
	SentAt  time.Time
}

func (n *Notification) GetIdentifier() string {
	switch n.Method {
	case enum.System:
		return n.UserID
	case enum.Email:
		return n.Email
	default:
		return ""
	}
}
