package controllers

import (
	"net/http"

	"github.com/H-Richard/talent/api/responses"
)

// Health function to ping postgres [WIP]
func (server *Server) Health(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK , "Server is up and running")
}