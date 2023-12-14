package server

import (
	"log"
	"notification-service/db/entity"
	"notification-service/db/query"
	"notification-service/db/repository"
	"notification-service/enum"
	"notification-service/mail"
	"notification-service/push_notification"
	"sync"
	"time"
)

var (
	muSend           sync.Mutex
	muSendBatch      sync.Mutex
	runningSend      bool
	runningSendBatch bool
)

type Watcher struct {
	notificationRepo       repository.NotificationRepository
	mail                   mail.Mailer
	system                 push_notification.Sender
	maxAmountNotifications int

	group     map[notificationGroup]compoundNotification
	processed map[int64]bool

	timeLimitReached chan struct{}
}

func NewWatcher(notificationRepo repository.NotificationRepository, mail mail.Mailer, system push_notification.Sender, maxAmountNotifications int) *Watcher {
	return &Watcher{
		notificationRepo:       notificationRepo,
		mail:                   mail,
		system:                 system,
		maxAmountNotifications: maxAmountNotifications,
		group:                  make(map[notificationGroup]compoundNotification),
		processed:              make(map[int64]bool),
		timeLimitReached:       make(chan struct{}),
	}
}

type notificationGroup struct {
	name       string
	method     enum.DeliveryMethod
	identifier string
}

type compoundNotification struct {
	compoundContent string
	notifications   []*entity.Notification
}

func (s *Watcher) Send() {
	muSend.Lock()
	defer muSend.Unlock()

	if runningSend {
		return
	}

	runningSend = true

	// get all notifications of type instantaneous not sent
	notifications, err := s.notificationRepo.Find(query.NotificationQuery{Sent: false, Type: enum.Instantaneous})
	if err != nil {
		log.Println("error getting instantaneous notifications")
	}

	s.sendNotifications(notifications)
	for _, notification := range notifications {
		notification.SentAt = time.Now()
		err = s.notificationRepo.Update(notification)
		if err != nil {
			log.Println("error updating notification status")
		}
	}

	runningSend = false
}

func (s *Watcher) SendBatch() {
	muSendBatch.Lock()
	defer muSendBatch.Unlock()

	if runningSendBatch {
		return
	}

	runningSendBatch = true

	// array of notifications to send
	var toSend []*entity.Notification

	select {
	// if time limit reached, then send all grouped notifications until now
	case <-s.timeLimitReached:
		for k, v := range s.group {
			n := new(entity.Notification)
			n.Name = k.name
			n.Method = k.method
			n.Content = v.compoundContent
			if k.method == enum.Email {
				n.Email = k.identifier
			}
			if k.method == enum.System {
				n.UserID = k.identifier
			}

			toSend = append(toSend, n)
		}

		s.group = make(map[notificationGroup]compoundNotification)
	default:
		// get all notifications of type batch not sent
		notifications, err := s.notificationRepo.Find(query.NotificationQuery{Sent: false, Type: enum.Batch})
		if err != nil {
			log.Println("error getting batch notifications")
		}

		// group notifications by name, method and identifier
		for _, notification := range notifications {
			if s.processed[notification.ID] {
				continue
			}

			s.processed[notification.ID] = true

			ng := notificationGroup{
				name:       notification.Name,
				method:     notification.Method,
				identifier: notification.GetIdentifier(),
			}

			cn := s.group[ng]
			cn.notifications = append(cn.notifications, notification)
			if cn.compoundContent != "" {
				cn.compoundContent += "\n"
			}
			cn.compoundContent += notification.Content
			s.group[ng] = cn

			// if max amount reached, then send notification
			if len(cn.notifications) >= s.maxAmountNotifications {
				nSend := notification
				nSend.Content = cn.compoundContent
				for _, n := range cn.notifications {
					delete(s.processed, n.ID)

					n.SentAt = time.Now()
					err = s.notificationRepo.Update(n)
					if err != nil {
						log.Println("error updating notification status")
					}
				}
				delete(s.group, ng)
				toSend = append(toSend, nSend)
			}
		}
	}

	s.sendNotifications(toSend)

	runningSendBatch = false
}

func (s *Watcher) TimeLimitReached() {
	s.timeLimitReached <- struct{}{}
}

func (s *Watcher) sendNotifications(notifications []*entity.Notification) {
	for _, n := range notifications {
		if n.Method == enum.Email {
			err := s.mail.Send(n.Email, n.Name, n.Content)
			if err != nil {
				log.Println("error sending notification by email")
				continue
			}
		}

		if n.Method == enum.System {
			err := s.system.Send(n.UserID, n.Name, n.Content)
			if err != nil {
				log.Println("error sending notification by system")
				continue
			}
		}
	}
}
