package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
)

type checkoutService struct {
	repo          *repository.CartRepository
	productRepo   *repository.ProductRepository
	inventoryRepo *repository.InventoryRepository
}

type CheckoutRequest struct {
	CartID string `json:"cart_id"`
	UserID string `json:"user_id"`
}

type CheckoutResponse struct {
	OrderID       string    `json:"order_id"`
	TotalPrice    float64   `json:"total_price"`
	NumberOfItems int       `json:"number_of_items"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *checkoutService) Checkout(ctx context.Context, req *connect.Request[CheckoutRequest]) (*connect.Response[CheckoutResponse], error) {
	// cart Get Cart
	cart, err := s.repo.GetCart(ctx, req.CartID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	// Lock inventory
	for _, product := range cart.Products {
		inventory, err := s.inventoryRepo.LockInventory(ctx, product.ProductID, product.Quantity)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	// Update Checkout status
	err = s.repo.UpdateCheckoutStatus(ctx, req.CartID, "locked")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// Process payment
	// Create Payment Clinet
	// get Succesfull Payment response form payment clint
	// Update Checkout status
	err = s.repo.UpdateCheckoutStatus(ctx, req.CartID, "paymnet_successful")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create Order
	order, err := s.repo.CreateOrder(ctx, req.CartID, req.UserID)
	if err != nil {
		// release inventory lock
		// update checkout status as failed()
		return nil, connect.NewError(connect.CodeInternal, err)

	}
	// release/ remove lock from inventory
	for _, product := range cart.Products {
		err = s.inventoryRepo.ReleaseInventory(ctx, product.ProductID, product.Quantity)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Update Checkout status
	err = s.repo.UpdateCheckoutStatus(ctx, req.CartID, "completed")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// empty cart
	err = s.repo.EmptyCart(ctx, req.CartID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&CheckoutResponse{
		OrderID:       order.ID,
		TotalPrice:    order.TotalPrice,
		NumberOfItems: order.NumberOfItems,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
	}), nil
}
