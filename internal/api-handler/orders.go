package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Find user Orders using status.
// (GET /store/orders)
func (a *APIHandler) FindOrders(ctx context.Context,
	request genRouter.FindOrdersRequestObject) (genRouter.FindOrdersResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Place orders for pets.
// (POST /store/orders)
func (a *APIHandler) PlaceOrders(ctx context.Context,
	request genRouter.PlaceOrdersRequestObject) (genRouter.PlaceOrdersResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete user order by identifier.
// (DELETE /store/orders/{orderId})
func (a *APIHandler) DeleteOrder(ctx context.Context,
	request genRouter.DeleteOrderRequestObject) (genRouter.DeleteOrderResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find user order by ID.
// (GET /store/orders/{orderId})
func (a *APIHandler) GetOrderById(ctx context.Context,
	request genRouter.GetOrderByIdRequestObject) (genRouter.GetOrderByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
