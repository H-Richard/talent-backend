package controllers

import (
	"github.com/H-Richard/talent/api/middleware"
)

var (
	// GET HTTP Keyword
	GET = "GET"
	// POST HTTP Keyword
	POST = "POST"
)

func (s *Server) initializeRoutes() {
	// GET /healthz 
	s.Router.HandleFunc("/healthz", middleware.SetMiddlewareJSON(s.Health)).Methods(GET)

	// POST /login
	s.Router.HandleFunc("/login", middleware.SetMiddlewareJSON(s.Login)).Methods(POST)

	// Users routes
	s.Router.HandleFunc("/users", middleware.SetMiddlewareJSON(s.CreateUser)).Methods(POST)
	s.Router.HandleFunc("/users", middleware.SetMiddlewareAuthentication(s.GetUsers)).Methods(GET)

	// Posts routes
	s.Router.HandleFunc("/posts", middleware.SetMiddlewareJSON(s.GetActivePosts)).Methods(GET)
	s.Router.HandleFunc("/posts", middleware.SetMiddlewareJSON(s.CreatePost)).Methods(POST)
}