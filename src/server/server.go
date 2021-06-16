package server

import (
	"api-note/src/model"
	"api-note/src/store"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	adress string
	router *mux.Router
	store  *store.Store
}

// Create new server
func NewServer(pattern string) *Server {
	return &Server{
		adress: pattern,
		router: mux.NewRouter(),
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
		log.Printf("_%v_", email)

		usr, _ := s.store.User().GetUserByEmail(email)
		log.Println(usr)
		data, _ := json.Marshal(usr)

		rw.Write(data)
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

		_, err := s.store.User().Create(usr)

		if err != nil {
			data := struct {
				Err error
			}{
				err,
			}
			res, err := json.Marshal(data)
			if err != nil {
				log.Println(err)
			}
			w.Write(res)
			return
		}
	}
}
