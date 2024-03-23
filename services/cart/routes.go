package cart

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/trenchesdeveloper/go-ecom/services/auth"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
	"net/http"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
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

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"orderID":    orderID,
		"totalPrice": totalPrice,
	})
}
