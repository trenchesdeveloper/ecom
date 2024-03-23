package product

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", h.GetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", h.CreateProduct).Methods(http.MethodPost)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	product, err := h.store.GetProduct(idInt)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductInput
	err := utils.ParseJSON(w, r, &product)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate input
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, fmt.Errorf("invalid payload %v", errors), http.StatusBadRequest)
		return

	}

	log.Println("Product: ", product)

	err = h.store.CreateProduct(types.Product{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Quantity:    product.Quantity,
	})

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}
