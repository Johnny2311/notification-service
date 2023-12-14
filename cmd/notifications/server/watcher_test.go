package server

import (
	"github.com/stretchr/testify/assert"
	"notification-service/db/entity"
	"notification-service/db/query"
	"notification-service/db/repository"
	"notification-service/enum"
	"notification-service/mail"
	"notification-service/push_notification"
	"testing"
)

func Test_WatcherSend(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	mailProvider := mail.NewYahooMailProvider()
	systemProvider := push_notification.NewSystemNotificationSender()

	watcher := NewWatcher(repo, mailProvider, systemProvider, 10)

	notification := &entity.Notification{Type: enum.Instantaneous}
	err := repo.Create(notification)
	assert.NoError(t, err)
	err = repo.Create(notification)
	assert.NoError(t, err)
	err = repo.Create(notification)
	assert.NoError(t, err)

	watcher.Send()

	found, err := repo.Find(query.NotificationQuery{Sent: false})
	assert.NoError(t, err)
	assert.Empty(t, found)
}

func Test_WatcherSendBatch(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	mailProvider := mail.NewYahooMailProvider()
	systemProvider := push_notification.NewSystemNotificationSender()

	watcher := NewWatcher(repo, mailProvider, systemProvider, 2)

	notification1 := &entity.Notification{ID: 1, Type: enum.Batch}
	notification2 := &entity.Notification{ID: 2, Type: enum.Batch}
	notification3 := &entity.Notification{ID: 3, Type: enum.Batch}
	err := repo.Create(notification1)
	assert.NoError(t, err)
	err = repo.Create(notification2)
	assert.NoError(t, err)
	err = repo.Create(notification3)
	assert.NoError(t, err)

	watcher.SendBatch()

	// max amount is 2, so must remain only 1
	found, err := repo.Find(query.NotificationQuery{Sent: false})
	assert.NoError(t, err)
	assert.Len(t, found, 1)
}
