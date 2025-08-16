package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Find user Orders using status.
// (GET /store/orders)
func (a *ApiHandler) FindOrders(ctx context.Context, request genRouter.FindOrdersRequestObject) (genRouter.FindOrdersResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Place orders for pets.
// (POST /store/orders)
func (a *ApiHandler) PlaceOrders(ctx context.Context, request genRouter.PlaceOrdersRequestObject) (genRouter.PlaceOrdersResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete purchase order by identifier.
// (DELETE /store/orders/{orderId})
func (a *ApiHandler) DeleteOrder(ctx context.Context, request genRouter.DeleteOrderRequestObject) (genRouter.DeleteOrderResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find purchase order by ID.
// (GET /store/orders/{orderId})
func (a *ApiHandler) GetOrderById(ctx context.Context, request genRouter.GetOrderByIdRequestObject) (genRouter.GetOrderByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
