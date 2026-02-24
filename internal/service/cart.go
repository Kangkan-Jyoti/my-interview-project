package service


type cartService struct {
	repo *repository.CartRepository

}

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UserID    string `json:"user_id"`
	CartID    string `json:"cart_id"`
}

func (s *cartService) AddToCart(ctx context.Context, req *connect.Request[AddToCartRequest]) (*connect.Response[AddToCartResponse], error) {

	if req.Quantity <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("quantity must be greater than 0"))
	}
	// Check product is exist 

	_, err := s.repo.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// Check inventory returning number of total_stock
	// row level locking 
	stock , err := s.repo.GetInventory(ctx, req.ProductID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if stock < req.Quantity {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not enough stock"))
	}

	// Add to cart
	// Need a vliadtion helper to valadet request parameters
	if req.UserID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}
	_, err = s.repo.AddToCart(ctx, req.UserID, req.ProductID, req.Quantity)
	// Check if the cart is already exist
	// Update / add to carts_product
	// add a new cart if not exist
	// Update Inventory with product_id and quantity 
	// Return the cart_id, user_id, product_id, quantity
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&AddToCartResponse{
		Message: "Product added to cart successfully",
	}), nil

}