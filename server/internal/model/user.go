package model

type User struct {
	ID          int    `json:"id,omitempty"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty"`
	// SystemUserName string    `json:"system_user_name"`
	// SystemUserId   string    `json:"system_user_id"`
}
