package query

import (
	"notification-service/enum"
)

type NotificationQuery struct {
	ID     int64
	Method enum.DeliveryMethod
	Type   enum.NotificationType
	UserID string
	Sent   bool
	Seen   bool
}
