package models

import fb "github.com/huandu/facebook"

var f *Facebook

func init() {
	f = new(Facebook)
}

type Facebook struct {
	PageID  string
	Session *fb.Session
}

func GetFacebook() *Facebook {
	return f
}
