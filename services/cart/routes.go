package cart

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
	"log"
	"net/http"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore) *Handler {
	return &Handler{store: store, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handleCheckout).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var checkoutInput types.CheckoutInput

	if err := utils.ParseJSON(w, r, &checkoutInput); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(checkoutInput); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, fmt.Errorf("invalid payload %v", errors), http.StatusBadRequest)
		return
	}

	// get products
	productIDs, err := getCartItemsIDs(checkoutInput.Items)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return

	}
	ps, err := h.productStore.GetProductByIDs(productIDs)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, checkoutInput.Items, userID)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	log.Printf("order created %d", orderID)
	log.Println("total price", totalPrice)

	utils.WriteJSON(w, http.StatusOK, nil)
}
