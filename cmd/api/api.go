package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/services/cart"
	"github.com/trenchesdeveloper/go-ecom/services/order"
	"github.com/trenchesdeveloper/go-ecom/services/product"
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

	// product service
	productStore := product.NewStore(a.db)

	productHandler := product.NewHandler(productStore)

	productHandler.RegisterRoutes(subrouter)

	// order service
	orderStore := order.NewStore(a.db)

	// cart service
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)

	cartHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on port", a.addr)

	return http.ListenAndServe(a.addr, router)
}
