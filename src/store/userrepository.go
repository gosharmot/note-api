package store

import (
	"api-note/src/model"
	"log"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) Create(usr model.User) (*model.User, error) {
	if err := u.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, username) VALUES ($1, $2, $3) RETURNING id",
		usr.Email,
		usr.EncryptedPassword,
		usr.Username,
	).Scan(&usr.ID); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (u *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	usr := &model.User{}
	if err := u.store.db.QueryRow(
		"SELECT id, email, encrypted_password, username FROM users WHERE email = $1",
		email,
	).Scan(
		&usr.ID,
		&usr.Email,
		&usr.EncryptedPassword,
		&usr.Username,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	return usr, nil
}
