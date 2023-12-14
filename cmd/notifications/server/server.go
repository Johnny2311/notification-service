package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"notification-service/db/repository"
	"notification-service/mail"
	"notification-service/push_notification"
	"os"
	"strconv"
)

func Run() {
	// creating notification repository
	notificationRepo := repository.NewInMemoryNotificationRepository()

	// creating notification handler
	h := NewNotificationHandler(notificationRepo)

	// configuring router
	r := GetRouter(h)

	amountStr := os.Getenv("NOTIFICATION_BATCH_AMOUNT")
	maxAmountNotifications, err := strconv.Atoi(amountStr)
	if err != nil {
		log.Fatal(err)
	}

	// configuring and starting cron
	mailProvider := mail.NewYahooMailProvider()
	systemProvider := push_notification.NewSystemNotificationSender()
	watcher := NewWatcher(notificationRepo, mailProvider, systemProvider, maxAmountNotifications)
	StartCron(watcher)

	// starting server
	address := os.Getenv("HOST")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

func GetRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	// configuring subrouter for api
	v1Router := r.PathPrefix("/api/v1").Subrouter()
	v1Router.HandleFunc("/notifications/send", h.SendNotificationHandler).Methods(http.MethodPost)
	v1Router.HandleFunc("/notifications/{user_id}/list", h.ListNotificationHandler).Methods(http.MethodGet)
	v1Router.HandleFunc("/notifications/remove/{id}", h.RemoveNotificationHandler).Methods(http.MethodDelete)
	v1Router.HandleFunc("/notifications/mark/{id}", h.MarkNotificationHandler).Methods(http.MethodGet)

	return r
}

func StartCron(watcher *Watcher) {
	batchTimeStr := os.Getenv("NOTIFICATION_BATCH_TIME_SEG")
	batchTime, err := strconv.ParseInt(batchTimeStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New()
	err = c.AddFunc("* * * * *", watcher.Send)
	if err != nil {
		log.Fatal(err)
	}
	err = c.AddFunc("* * * * *", watcher.SendBatch)
	if err != nil {
		log.Fatal(err)
	}
	err = c.AddFunc(fmt.Sprintf("*/%d * * * *", batchTime), watcher.TimeLimitReached)
	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
