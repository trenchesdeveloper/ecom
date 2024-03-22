package user

import (
	"github.com/gorilla/mux"
	go_ecom "github.com/trenchesdeveloper/go-ecom"
	"github.com/trenchesdeveloper/go-ecom/services/auth"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
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
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
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
	if err := payload.Validate(); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	//check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.ErrorJSON(w, go_ecom.ErrUserExists, http.StatusBadRequest)
		return
	}

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

	utils.WriteJSON(w, http.StatusCreated, nil)
}
