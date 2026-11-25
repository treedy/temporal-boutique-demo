// package "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/tmprl"
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

	var orderItems []*pb.CartItem
	err = workflow.ExecuteActivity(ctx, GetUserCart, req.UserId).Get(ctx, &orderItems)
	if err != nil {
		logger.Error("Failed to get cart", "for user", req.UserId)
		return false, err
	} else {
		if len(orderItems) < 1 {
			logger.Error("No items in cart", "ID", req.UserId)
		}
		for _, item := range orderItems {
			logger.Info("Item in cart", "Item:", item.String())
		}
	}

	logger.Debug("Continuing for order", "ID", orderID)

	return true, nil
}
