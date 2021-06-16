package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	databaseUrl    string
	db             *sql.DB
	userRepository *UserRepository
}

func NewStore() *Store {
	return &Store{
		databaseUrl: "host=localhost dbname=restapi_note sslmode=disable",
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.databaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
