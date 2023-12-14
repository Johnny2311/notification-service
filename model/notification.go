package model

import (
	"notification-service/db/entity"
	"notification-service/enum"
	"time"
)

type Notification struct {
	Date           time.Time             `json:"date"`
	Name           string                `json:"name"`
	Method         enum.DeliveryMethod   `json:"method"`
	Type           enum.NotificationType `json:"type"`
	UserIdentifier string                `json:"user_identifier"`
	Content        string                `json:"content"`
}

func (n *Notification) ToEntity() *entity.Notification {
	notification := new(entity.Notification)
	notification.Date = n.Date
	notification.Name = n.Name
	notification.Method = n.Method
	notification.Type = n.Type
	notification.Content = n.Content

	switch n.Method {
	case enum.Email:
		notification.Email = n.UserIdentifier
	case enum.System:
		notification.UserID = n.UserIdentifier
	}

	return notification
}
