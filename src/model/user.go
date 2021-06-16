package model

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
	Username          string `json:"username"`
}
