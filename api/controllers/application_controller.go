package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/H-Richard/talent/api/auth"
	"github.com/H-Richard/talent/api/models"
	"github.com/H-Richard/talent/api/responses"
	"github.com/gorilla/mux"
)

// CreateApplication creates applications
func (server *Server) CreateApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	fmt.Println(postID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	application := models.Application{}
	application.PostID = uint32(postID)
	err = json.Unmarshal(body, &application)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	application.Pre()
	err = application.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Token Unauthorized"))
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Token Unauthorized"))
		return
	}
	user := models.User{}
	foundUser, err := user.FindByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if foundUser.Executive {
		responses.ERROR(w, http.StatusTeapot, errors.New("Executives should not be applying for jobs externally"))
		return
	}
	application.ApplicantID = uid
	post := models.Post{}
	foundPost, err := post.FindByID(server.DB, application.PostID)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if !foundPost.Active {
		responses.ERROR(w, http.StatusGone, err)
		return
	}
	createdApplication, err := application.SaveApplication(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, createdApplication.ID))
	responses.JSON(w, http.StatusCreated, createdApplication.JSON())
}
