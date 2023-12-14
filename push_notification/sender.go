package push_notification

type Sender interface {
	Send(userID, eventName string, data any) error
}
