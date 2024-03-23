package cart

import (
	"fmt"
	"github.com/trenchesdeveloper/go-ecom/types"
)

func getCartItemsIDs(cartItems []types.CartItem) ([]int, error) {
	productIds := make([]int, len(cartItems))

	for i, cartItem := range cartItems {
		if cartItem.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", cartItem.ProductID)
		}
		productIds[i] = cartItem.ProductID
	}

	return productIds, nil

}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	// check if all products are in stock
	for _, product := range products {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(productMap, items); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice := calculateTotalPrice(productMap, items)

	// reduce quantity of products
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		// update product
		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, err
		}

	}

	// create order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St", //TODO get address from user
	})

	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		product := productMap[item.ProductID]
		if err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}); err != nil {
			return 0, 0, err
		}

	}

	return 0, totalPrice, nil
}

func checkIfCartIsInStock(productMap map[int]types.Product, items []types.CartItem) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d not found", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is out of stock", item.ProductID)
		}
	}
	return nil
}

func calculateTotalPrice(productMap map[int]types.Product, items []types.CartItem) float64 {
	var totalPrice float64
	for _, item := range items {
		product := productMap[item.ProductID]
		totalPrice += float64(item.Quantity) * product.Price
	}
	return totalPrice
}
