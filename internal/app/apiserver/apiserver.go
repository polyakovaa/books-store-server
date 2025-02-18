package apiserver

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/polyakovaa/standartserver3/store"
	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api/v1"
)

// APIServer object for instancing server
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// constructor APIServer
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// configures loggers, router, db connection and starts server
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("starting api server at port :", s.config.BindAddr)
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return nil
	}
	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc(prefix+"/books", s.GetAllBooks).Methods("GET")
	s.router.HandleFunc(prefix+"/books"+"/{id}", s.GetBookById).Methods("GET")
	s.router.HandleFunc(prefix+"/books"+"/{id}", s.DeleteBookById).Methods("DELETE")
	s.router.HandleFunc(prefix+"/books", s.PostBook).Methods("POST")
	s.router.HandleFunc(prefix+"/user/register", s.PostUserRegister).Methods("POST")

}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}
