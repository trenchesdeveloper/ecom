package user

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/config"
	"github.com/trenchesdeveloper/go-ecom/services/auth"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
	"log"
	"net/http"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// parse json payload
	var payload types.LoginInput

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// sanitize input
	payload.Sanitize()

	// validate input
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, fmt.Errorf("invalid payload %v", errors), http.StatusBadRequest)
		return

	}

	// get user by email
	user, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("invalid email or password"), http.StatusBadRequest)
		return
	}

	// compare password
	if !auth.ComparePassword(user.Password, payload.Password) {
		utils.ErrorJSON(w, fmt.Errorf("invalid email or password"), http.StatusBadRequest)
		return
	}

	// generate token
	token, err := auth.GenerateToken(config.Envs.JWTSecret, user.ID)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return

	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// parse JSON payload
	var payload types.RegisterInput

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// sanitize input
	payload.Sanitize()

	// validate input
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, fmt.Errorf("invalid payload %v", errors), http.StatusBadRequest)
		return
	}

	//check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.ErrorJSON(w, fmt.Errorf("user already exists"), http.StatusBadRequest)
		return
	}

	log.Println("User does not exist")

	// hash password
	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return

	}

	// create user
	if err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  string(hashedPassword),
	}); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("User created successfully")

	utils.WriteJSON(w, http.StatusCreated, nil)

}
