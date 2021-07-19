package server

import (
	"api-note/src/model"
	"api-note/src/store"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Server struct {
	adress       string
	router       *mux.Router
	store        *store.Store
	sessionStore sessions.Store
}

// Create new server
func NewServer(pattern string, sessionStore sessions.Store) *Server {
	return &Server{
		adress:       pattern,
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
	}
}

func (s *Server) Listen() error {
	log.Println("Listening server...")

	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.adress, s.router)
}

//Configurate router
func (s *Server) configureRouter() {
	s.router.HandleFunc("/create/user/", s.createUser()).Methods("POST")
	s.router.HandleFunc("/user/{email}/", s.getUser()).Methods("GET")
	s.router.HandleFunc("/sessions/", s.createSession()).Methods("POST")
}

//Configurate store
func (s *Server) configureStore() error {
	store := store.NewStore()

	if err := store.Open(); err != nil {
		return err
	}

	s.store = store
	return nil
}

func (s *Server) getUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		email := mux.Vars(r)["email"]

		usr, _ := s.store.User().GetUserByEmail(email)
		if usr != nil {

			data := model.User{
				Email:    usr.Email,
				ID:       usr.ID,
				Username: usr.Username,
			}

			json.NewEncoder(rw).Encode(data)
		}
	}
}

func (s *Server) createUser() http.HandlerFunc {
	usr := model.User{}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&usr); err != nil {
			log.Println(err)
		}

		u, err := s.store.User().Create(usr)

		if err != nil {
			data := struct {
				Err error
			}{
				err,
			}
			res, _ := json.Marshal(data)

			w.Write(res)
			return
		}

		data := model.User{
			Email:    u.Email,
			ID:       u.ID,
			Username: u.Username,
		}

		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) createSession() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return
		}

		u, err := s.store.User().GetUserByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			return
		}

		session, _ := s.sessionStore.Get(r, "api_note")

		session.Values["user_id"] = u.ID

		s.sessionStore.Save(r, rw, session)
	}
}
