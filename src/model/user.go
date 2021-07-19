package model

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	Username          string `json:"username"`
}

func (u *User) ComparePassword(password string) bool {
	return u.EncryptedPassword == password
}
