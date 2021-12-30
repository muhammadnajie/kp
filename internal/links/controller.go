package links

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/muhammadnajie/kp/internal/auth"
	"github.com/muhammadnajie/kp/internal/resources"
	"github.com/muhammadnajie/kp/internal/users"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetAllController(w http.ResponseWriter, r *http.Request) {
	user := auth.ExtractUserContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	var dbLinks []Link
	var err error
	dbLinks, err = GetAll(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := resources.Success{
		Code: http.StatusOK,
		Data: dbLinks,
	}
	fmt.Println(res)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetByTitleController(w http.ResponseWriter, r *http.Request) {
	user := auth.ExtractUserContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	type T = struct {
		Title string `json:"title"`
	}
	var title T

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	err = json.Unmarshal(b, &title)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	var dbLinks []Link
	dbLinks, err = GetByTitle(title.Title, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := resources.Success{
		Code: http.StatusOK,
		Data: dbLinks,
	}
	fmt.Println(res)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func CreateLinkController(w http.ResponseWriter, r *http.Request) {
	user := auth.ExtractUserContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	var payload PayloadCreateLink
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	err = json.Unmarshal(b, &payload)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	var link Link
	link.Title = payload.Title
	link.Address = payload.Address
	link.User = &users.User{
		ID: user.ID,
	}
	linkID := link.Save()

	res := resources.Success{
		Code: http.StatusOK,
		Data: linkID,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UpdateLinkController(w http.ResponseWriter, r *http.Request) {
	user := auth.ExtractUserContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	var payload PayloadUpdateLink
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	err = json.Unmarshal(b, &payload)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	ID, _ := strconv.Atoi(payload.ID)
	userID, _ := strconv.Atoi(user.ID)
	link, err := GetByID(ID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resources.Failed{
				Code:    http.StatusNotFound,
				Message: "not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	link.Title = payload.Title
	link.Address = payload.Address
	_, err = link.Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res := resources.Success{
		Code:    http.StatusOK,
		Data:    nil,
		Message: "successfully update data",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func DeleteLinkController(w http.ResponseWriter, r *http.Request) {
	user := auth.ExtractUserContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	var payload PayloadDeleteLink
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	err = json.Unmarshal(b, &payload)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	var link Link
	link.ID = payload.ID
	_, err = link.Delete()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resources.Failed{
				Code:    http.StatusNotFound,
				Message: "not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resources.Failed{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res := resources.Success{
		Code:    http.StatusOK,
		Data:    nil,
		Message: "successfully delete data",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
