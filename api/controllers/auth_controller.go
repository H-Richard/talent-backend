package controllers

import (
	"encoding/json"
	"net/http"
	"io/ioutil"

	"github.com/H-Richard/talent/api/models"
	"github.com/H-Richard/talent/api/responses"
	"github.com/H-Richard/talent/api/auth"
	"github.com/H-Richard/talent/api/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignIn signs in the user
func (server *Server) SignIn(email, password string) (string ,error) {
	var err error
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.Create(user.ID)
}

// Login logs in the user
func (server *Server) Login (w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user.Pre()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, utils.FormatError(err.Error()))
		return
	}
	responses.JSON(w, http.StatusOK, token)
}
