package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/oxfeeefeee/appgo"
	"github.com/oxfeeefeee/appgo/auth"
	"github.com/rs/cors"
	"net/http"
	"strings"
)

type Server struct {
	ts TokenStore
	*mux.Router
}
type TokenStore interface {
	Validate(token auth.Token) bool
}

func NewServer(ts TokenStore) *Server {
	return &Server{
		ts,
		mux.NewRouter(),
	}
}

func (s *Server) AddRest(path string, rest []interface{}) {
	for _, api := range rest {
		h := newHandler(api, s.ts)
		s.Handle(path+h.path, h).Methods(h.supports...)
	}
}

func (s *Server) AddHandler(path string,
	f func(w http.ResponseWriter, s *http.Request)) {
	s.HandleFunc(path, f)
}

func (s *Server) Serve() {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewCustomMiddleware(
		appgo.Conf.LogLevel, &log.TextFormatter{}, "appgo"))
	n.Use(cors.New(corsOptions()))
	n.UseHandler(s)
	n.Run(appgo.Conf.Negroni.Port)
}

func corsOptions() cors.Options {
	origins := strings.Split(appgo.Conf.Cors.AllowedOrigins, ",")
	methods := strings.Split(appgo.Conf.Cors.AllowedMethods, ",")
	headers := strings.Split(appgo.Conf.Cors.AllowedHeaders, ",")
	return cors.Options{
		AllowedOrigins:     origins,
		AllowedMethods:     methods,
		AllowedHeaders:     headers,
		OptionsPassthrough: appgo.Conf.Cors.OptionsPassthrough,
		Debug:              appgo.Conf.Cors.Debug,
	}
}
