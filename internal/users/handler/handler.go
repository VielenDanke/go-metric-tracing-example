package handler

import (
	"encoding/json"
	"fmt"
	"github.com/vielendanke/opentracing-example/internal/users/model"
	"github.com/vielendanke/opentracing-example/internal/users/service"
	"log"
	"net/http"
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
}
