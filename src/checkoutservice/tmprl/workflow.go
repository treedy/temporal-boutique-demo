package tmprl

import (
	// "uuid"

	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
	"github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/money"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func PlaceOrderWorkflow(ctx workflow.Context, req *pb.PlaceOrderRequest) (bool, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	orderID, err := uuid.NewUUID()
	if err != nil {
		logger.Error("failed to generate order uuid", "error:", err.Error())
		return false, err
	}

	// GetUserCart
	var cartItems []*pb.CartItem
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).GetUserCart, req.UserId).Get(ctx, &cartItems)
	if err != nil {
		logger.Error("Failed to get cart", "for user", req.UserId)
		return false, err
	} else { // Remove else when implementing next steps
		if len(cartItems) < 1 {
			logger.Error("No items in cart", "ID", req.UserId)
		}
	}

	// PrepOrderItems
	var orderItems []*pb.OrderItem
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).PrepOrderItems, cartItems, req.UserCurrency).Get(ctx, &orderItems)
	if err != nil {
		logger.Error("Failed to prep items", "for order", orderID)
		return false, err
	} else { // Remove else when implementing next steps
		if len(orderItems) < 1 {
			logger.Error("No items in order", "ID", orderID)
		}
		for _, pricedItem := range orderItems {
			logger.Debug("", "Priced Item", pricedItem.String())

		}
	}

	// QuoteShipping
	var shippingCostLocalized *pb.Money
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).QuoteShipping, req.Address, cartItems, req.UserCurrency).Get(ctx, &shippingCostLocalized)
	if err != nil {
		logger.Error("Failed to get shipping costs")
		return false, err
	} else {
		logger.Debug("Shipping", "Calcualted Cost", shippingCostLocalized)
	}

	// Calculate Total
	total := pb.Money{CurrencyCode: req.UserCurrency,
		Units: 0,
		Nanos: 0}
	total = money.Must(money.Sum(total, *shippingCostLocalized))
	for _, it := range orderItems {
		multPrice := money.MultiplySlow(*it.Cost, uint32(it.GetItem().GetQuantity()))
		total = money.Must(money.Sum(total, multPrice))
	}
	// TODO(treedy): Return to customer for verificaiton

	// Charge card (single step, no product reservation)
	var txID string
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).ChargeCard, total, req.CreditCard).Get(ctx, &txID)
	if err != nil {
		// return false, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
		logger.Error("Failed to charge card")
		return false, err
	}
	logger.Info("payment went through", "transaction ID", txID)

	// Ship the order
	var shippingTrackingID string
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).ShipOrder, req.Address, cartItems).Get(ctx, &shippingTrackingID)

	// Finalize the order
	// Clears the user's cart and sends an email
	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       shippingCostLocalized,
		ShippingAddress:    req.Address,
		Items:              orderItems,
	}

	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).FinalizeOrder, req, orderResult).Get(ctx, nil)

	// _ = &pb.PlaceOrderResponse{Order: orderResult}
	// return resp, nil
	return true, nil
}
