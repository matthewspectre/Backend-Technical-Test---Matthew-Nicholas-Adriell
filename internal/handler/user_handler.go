package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"be_evindo/internal/entity/user"
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(usecase user.Usecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (handler *UserHandler) GetByID(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(strings.TrimPrefix(request.URL.Path, "/users/"), 10, 64)
	if err != nil {
		http.Error(response, "invalid user id", http.StatusBadRequest)
		return
	}

	result, err := handler.usecase.FindByID(request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, strconv.ErrSyntax) {
			status = http.StatusBadRequest
		}
		http.Error(response, err.Error(), status)
		return
	}
	writeJSON(response, http.StatusOK, result)
}

func (handler *UserHandler) Create(response http.ResponseWriter, request *http.Request) {
	var input user.User
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(response, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := handler.usecase.Create(request.Context(), &input); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(response, http.StatusCreated, input)
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}
