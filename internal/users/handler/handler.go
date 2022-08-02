package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vielendanke/opentracing-example/internal/users/model"
	"github.com/vielendanke/opentracing-example/internal/users/service"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return UserHandler{userService: service}
}

func (uh UserHandler) FindAll(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	usrs, err := uh.userService.FindAll(ctx)

	if err != nil {
		log.Println(err)
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, dataErr := json.Marshal(usrs)

	if dataErr != nil {
		log.Println(dataErr)
		rw.Write([]byte(dataErr.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(data)
}

func (uh UserHandler) Save(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user model.User

	if decErr := json.NewDecoder(r.Body).Decode(&user); decErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := uh.userService.Save(ctx, user)

	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(fmt.Sprintf("User saved: ID %d", id)))
	rw.WriteHeader(http.StatusCreated)
}

func (uh UserHandler) FindByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, findErr := uh.userService.FindByID(ctx, id)

	if findErr != nil {
		rw.Write([]byte(findErr.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, dataErr := json.Marshal(usr)

	if dataErr != nil {
		rw.Write([]byte(dataErr.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(data)
}

func (uh UserHandler) Update(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user model.User

	if decErr := json.NewDecoder(r.Body).Decode(&user); decErr != nil {
		rw.Write([]byte(decErr.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := uh.userService.Update(ctx, user); err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
