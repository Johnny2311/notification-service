package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"notification-service/db/entity"
	"notification-service/db/query"
	"notification-service/db/repository"
	"notification-service/enum"
	"notification-service/model"
	"testing"
	"time"
)

func Test_NotificationSend(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	h := NewNotificationHandler(repo)
	router := GetRouter(h)
	ts := httptest.NewServer(router)
	defer ts.Close()

	notification := model.Notification{
		Date:           time.Now(),
		Name:           "event_created",
		Method:         enum.System,
		Type:           enum.Instantaneous,
		UserIdentifier: "23",
		Content:        "some-content",
	}

	jsonData, err := json.Marshal(notification)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/notifications/send", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")

	notifications, err := repo.Find(query.NotificationQuery{})
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
}

func Test_NotificationList(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	h := NewNotificationHandler(repo)
	router := GetRouter(h)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// creating 3 notifications, 2 of the testing user
	userNotification := &entity.Notification{UserID: "1"}
	noUserNotification := &entity.Notification{UserID: "2"}
	err := h.notificationRepo.Create(userNotification)
	assert.NoError(t, err)
	err = h.notificationRepo.Create(userNotification)
	assert.NoError(t, err)
	err = h.notificationRepo.Create(noUserNotification)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/notifications/1/list", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")

	notifications, err := repo.Find(query.NotificationQuery{UserID: "1"})
	assert.NoError(t, err)
	assert.Len(t, notifications, 2)
}

func Test_NotificationMark(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	h := NewNotificationHandler(repo)
	router := GetRouter(h)
	ts := httptest.NewServer(router)
	defer ts.Close()

	notification := &entity.Notification{ID: 1}
	err := h.notificationRepo.Create(notification)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/notifications/mark/1", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")

	n, err := repo.FindOne(1)
	assert.NoError(t, err)
	assert.Equal(t, true, n.Seen)
}

func Test_NotificationDelete(t *testing.T) {
	repo := repository.NewInMemoryNotificationRepository()
	h := NewNotificationHandler(repo)
	router := GetRouter(h)
	ts := httptest.NewServer(router)
	defer ts.Close()

	notification := &entity.Notification{ID: 1}
	err := h.notificationRepo.Create(notification)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/notifications/remove/1", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")

	notifications, err := repo.Find(query.NotificationQuery{})
	assert.NoError(t, err)
	assert.Empty(t, notifications)
}
