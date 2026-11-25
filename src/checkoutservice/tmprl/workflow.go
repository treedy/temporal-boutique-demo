package tmprl

import (
	// "uuid"

	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
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

	logger.Debug("Continuing for order", "ID", orderID)

	var orderItems []*pb.OrderItem
	err = workflow.ExecuteActivity(ctx, (*CheckoutActivities).PrepOrderItems, cartItems).Get(ctx, &orderItems)
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

	return true, nil
}
