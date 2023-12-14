package enum

type NotificationType string

const (
	Instantaneous NotificationType = "instantaneous"
	Batch         NotificationType = "batch"
)
