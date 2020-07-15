package controllers 

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/H-Richard/talent/api/responses"
	"github.com/H-Richard/talent/api/models"
	"github.com/H-Richard/talent/api/utils"
)

// CreateUser creates new users
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	user.Pre()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	createdUser, err := user.SaveUser(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, utils.FormatError(err.Error()))
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	responses.JSON(w, http.StatusCreated, createdUser)
}

// GetUser responds with a user if found
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	gottenUser, err := user.FindByID(server.DB, uint32(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, gottenUser)
}

// GetUsers responds with all users
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.All(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}