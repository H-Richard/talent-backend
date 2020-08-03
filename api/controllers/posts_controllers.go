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
	"github.com/H-Richard/talent/api/utils"
	"github.com/gorilla/mux"
)

// CreatePost creats posts
func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post.Pre()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
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
	if !foundUser.Executive {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	post.AuthorID = uid
	createdPost, err := post.SavePost(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, createdPost.ID))
	responses.JSON(w, http.StatusCreated, createdPost.JSON())
}

// GetActivePosts gets all current active posts
func (server *Server) GetActivePosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllActivePosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	var res []map[string]interface{} = []map[string]interface{}{}
	for i := range *posts {
		res = append(res, (*posts)[i].JSON())
	}
	responses.JSON(w, http.StatusOK, res)
}

// UpdatePost checks if the user is authorized and updates the post
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}
	// foundPost, err := post.FindByID(server.DB, uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post Not Found"))
		return
	}
	user := models.User{}
	foundUser, err := user.FindByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if !foundUser.Executive {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	postUpdate.Pre()
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	postUpdate.ID = post.ID
	postUpdated, err := postUpdate.UpdatePost(server.DB)
	if err != nil {
		FormatError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, FormatError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated.JSON())
}

// DeletePost delets a post only if the user is the author
func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = post.DeletePost(server.DB, uint32(pid), uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
