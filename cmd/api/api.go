package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/services/user"
	"log"
	"net/http"
)

type Application struct {
	addr string
	db   *sql.DB
}

func NewApplication(addr string, db *sql.DB) *Application {
	return &Application{
		addr: addr,
		db:   db,
	}
}

func (a *Application) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(a.db)

	userHandler := user.NewHandler(userStore)

	userHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on port", a.addr)

	return http.ListenAndServe(a.addr, router)
}
