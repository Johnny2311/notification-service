package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"notification-service/db/query"
	"notification-service/db/repository"
	"notification-service/model"
	"strconv"
)

type Handler struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationHandler(notificationRepo repository.NotificationRepository) *Handler {
	return &Handler{notificationRepo: notificationRepo}
}

func (h *Handler) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notification model.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid format in JSON", http.StatusBadRequest)
		return
	}

	if err = h.notificationRepo.Create(notification.ToEntity()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Notification will be sent"))
}

func (h *Handler) ListNotificationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID := params["user_id"]

	notifications, err := h.notificationRepo.Find(query.NotificationQuery{UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(notifications)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) MarkNotificationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationRepo.FindOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification.Seen = !notification.Seen
	if err = h.notificationRepo.Update(notification); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Notification marked as seen"))
}

func (h *Handler) RemoveNotificationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err = h.notificationRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Notification deleted"))
}
