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
		// return nil, status.Errorf(codes.Internal, "failed to generate order uuid")
		logger.Error("failed to generate order uuid", "error:", err.Error())
		return false, err
	}

	var orderItems []*pb.CartItem
	err = workflow.ExecuteActivity(ctx, GetUserCart, req.UserId).Get(ctx, &orderItems)
	if err != nil {
		logger.Error("Failed to get cart")
		return false, err
	} else {
		for _, item := range orderItems {
			logger.Info("Item in cart", "Item:", item.String())
		}
	}

	return true, nil
}
