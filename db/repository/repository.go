package repository

import (
	"notification-service/db/entity"
	"notification-service/db/query"
)

type NotificationRepository interface {
	Create(*entity.Notification) error
	Update(*entity.Notification) error
	Find(query.NotificationQuery) ([]*entity.Notification, error)
	FindOne(int64) (*entity.Notification, error)
	Delete(int64) error
}
