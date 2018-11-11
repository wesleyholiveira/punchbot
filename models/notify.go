package models

type Notify struct {
	UserID   string
	Projects *[]Project
	VIP      bool
}

var notifyUser TNotify

func init() {
	notifyUser = make(TNotify)
}

func GetNotifyUser() TNotify {
	return notifyUser
}
