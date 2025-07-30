package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Returns pet inventories by status.
// (GET /store/inventory)
func (h *apiHandler) GetInventory(ctx context.Context, request genRouter.GetInventoryRequestObject) (genRouter.GetInventoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Place an order for a pet.
// (POST /store/order)
func (h *apiHandler) PlaceOrder(ctx context.Context, request genRouter.PlaceOrderRequestObject) (genRouter.PlaceOrderResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete purchase order by identifier.
// (DELETE /store/order/{orderId})
func (h *apiHandler) DeleteOrder(ctx context.Context, request genRouter.DeleteOrderRequestObject) (genRouter.DeleteOrderResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find purchase order by ID.
// (GET /store/order/{orderId})
func (h *apiHandler) GetOrderById(ctx context.Context, request genRouter.GetOrderByIdRequestObject) (genRouter.GetOrderByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
