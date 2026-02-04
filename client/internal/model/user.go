package model

type User struct {
	UserName       string `json:"user_name,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Token          string `json:"token,omitempty"`
	StatusError    error  `json:"-"`
	SystemUserName string `json:"system_user_name"`
	SystemUserId   string `json:"system_user_id"`
}

func NewUser(userName, email, phone, password string) *User {
	return &User{
		UserName:    userName,
		Email:       email,
		PhoneNumber: phone,
		Password:    password,
	}
}
