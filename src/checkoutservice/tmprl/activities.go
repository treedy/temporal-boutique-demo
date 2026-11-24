// package github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/tmprl
package tmprl

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// ACTIVITY: GetUserCart
func GetUserCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	var cartSvcAddr string
	mustMapEnv(&cartSvcAddr, "CART_SERVICE_ADDR")
	var cartSvcConn *grpc.ClientConn
	mustConnGRPC(ctx, &cartSvcConn, cartSvcAddr)

	cart, err := pb.NewCartServiceClient(cartSvcConn).GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
	}
	return cart.GetItems(), nil
}

// GRPC Related functions
func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
	var err error
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	*conn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure())
	// grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	// grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	// *conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}
