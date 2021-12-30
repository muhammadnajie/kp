package users

import (
	"encoding/json"
	"errors"
	"github.com/muhammadnajie/kp/internal/pkg/jwt"
	"github.com/muhammadnajie/kp/internal/resources"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	valid, err := resources.ValidateUser(user.Username, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	authenticated := user.Authenticate()
	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: errors.New("invalid username or password").Error(),
		})
		return
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := resources.Success{
		Code: http.StatusOK,
		Data: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
	}
	json.NewEncoder(w).Encode(res)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	valid, err := resources.ValidateUser(user.Username, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := resources.Success{
		Code: http.StatusOK,
		Data: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
	}
	json.NewEncoder(w).Encode(res)
}
