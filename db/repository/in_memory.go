package repository

import (
	"errors"
	"notification-service/db/entity"
	"notification-service/db/query"
)

type notification struct {
	items []*entity.Notification
}

func NewInMemoryNotificationRepository() NotificationRepository {
	return &notification{}
}

func (n *notification) Create(e *entity.Notification) error {
	n.items = append(n.items, e)

	return nil
}

func (n *notification) Update(e *entity.Notification) error {
	for idx, item := range n.items {
		if item.ID == e.ID {
			n.items[idx] = e
		}
	}

	return nil
}

func (n *notification) Find(q query.NotificationQuery) ([]*entity.Notification, error) {
	var ret []*entity.Notification
	for _, item := range n.items {
		if q.ID != 0 {
			if item.ID != q.ID {
				continue
			}
		}
		if q.UserID != "" {
			if item.UserID != q.UserID {
				continue
			}
		}
		if q.Type != "" {
			if item.Type != q.Type {
				continue
			}
		}
		if q.Method != "" {
			if item.Method != q.Method {
				continue
			}
		}
		switch q.Sent {
		case true:
			if item.SentAt.IsZero() {
				continue
			}
		case false:
			if !item.SentAt.IsZero() {
				continue
			}
		}

		ret = append(ret, item)
	}

	return ret, nil
}

func (n *notification) FindOne(id int64) (*entity.Notification, error) {
	for _, item := range n.items {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, errors.New("not found")
}

func (n *notification) Delete(id int64) error {
	var nNew []*entity.Notification
	for idx, item := range n.items {
		if item.ID == id {
			nNew = append(nNew, n.items[:idx]...)

			if idx+1 < len(n.items) {
				nNew = append(nNew, n.items[idx+1:]...)
			}

			n.items = nNew
			return nil
		}
	}

	return errors.New("not found")
}
