package model

type User struct {
	UserName       string `json:"user_name,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	SystemUserName string `json:"system_user_name"`
	SystemUserId   string `json:"system_user_id"`
}
