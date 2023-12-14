package push_notification

type SystemNotificationSender struct{}

func NewSystemNotificationSender() Sender {
	return &SystemNotificationSender{}
}

func (s *SystemNotificationSender) Send(userID, eventName string, data any) error {
	// mock sending system notifications
	return nil
}
