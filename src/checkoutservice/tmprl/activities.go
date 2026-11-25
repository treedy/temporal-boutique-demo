// package github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/tmprl
package tmprl

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"

	"github.com/pkg/errors"
	"go.temporal.io/sdk/log"
	"google.golang.org/grpc"
)

// moved from main unedited (except for name)
type CheckoutActivities struct {
	pb.UnimplementedCheckoutServiceServer

	productCatalogSvcAddr string
	productCatalogSvcConn *grpc.ClientConn

	cartSvcAddr string
	cartSvcConn *grpc.ClientConn

	currencySvcAddr string
	currencySvcConn *grpc.ClientConn

	shippingSvcAddr string
	shippingSvcConn *grpc.ClientConn

	emailSvcAddr string
	emailSvcConn *grpc.ClientConn

	paymentSvcAddr string
	paymentSvcConn *grpc.ClientConn

	logger log.Logger
}

// ACTIVITY: GetUserCart
func (ca *CheckoutActivities) GetUserCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	cart, err := pb.NewCartServiceClient(ca.cartSvcConn).GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
	}
	return cart.GetItems(), nil
}

func (ca *CheckoutActivities) PrepOrderItems(ctx context.Context, items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
	out := make([]*pb.OrderItem, len(items))
	cl := pb.NewProductCatalogServiceClient(ca.productCatalogSvcConn)

	for i, item := range items {
		product, err := cl.GetProduct(ctx, &pb.GetProductRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, fmt.Errorf("failed to get product #%q", item.GetProductId())
		}
		price, err := ca.convertCurrency(ctx, product.GetPriceUsd(), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{
			Item: item,
			Cost: price}
	}
	return out, nil
}

// Currency converter is currently broken?
func (ca *CheckoutActivities) convertCurrency(ctx context.Context, from *pb.Money, toCurrency string) (*pb.Money, error) {
	return from, nil
	/*
		ca.logger.Debug("Converting", "from", from, "to", toCurrency)
		result, err := pb.NewCurrencyServiceClient(ca.currencySvcConn).Convert(ctx, &pb.CurrencyConversionRequest{
			From:   from,
			ToCode: toCurrency})
		if err != nil {
			return nil, fmt.Errorf("failed to convert currency: %+v", err)
		}
		ca.logger.Debug("Conversion done.", "result", result)
		return result, err
	*/
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

func NewCheckoutActivies(ctx context.Context, logger log.Logger) *CheckoutActivities {
	ca := &CheckoutActivities{}
	ca.logger = logger

	mustMapEnv(&ca.productCatalogSvcAddr, "PRODUCT_CATALOG_SERVICE_ADDR")
	mustMapEnv(&ca.cartSvcAddr, "CART_SERVICE_ADDR")
	mustMapEnv(&ca.currencySvcAddr, "CURRENCY_SERVICE_ADDR")
	mustMapEnv(&ca.shippingSvcAddr, "SHIPPING_SERVICE_ADDR")
	mustMapEnv(&ca.emailSvcAddr, "EMAIL_SERVICE_ADDR")
	mustMapEnv(&ca.paymentSvcAddr, "PAYMENT_SERVICE_ADDR")

	mustConnGRPC(ctx, &ca.productCatalogSvcConn, ca.productCatalogSvcAddr)
	mustConnGRPC(ctx, &ca.cartSvcConn, ca.cartSvcAddr)
	mustConnGRPC(ctx, &ca.currencySvcConn, ca.currencySvcAddr)
	mustConnGRPC(ctx, &ca.shippingSvcConn, ca.shippingSvcAddr)
	mustConnGRPC(ctx, &ca.emailSvcConn, ca.emailSvcAddr)
	mustConnGRPC(ctx, &ca.paymentSvcConn, ca.paymentSvcAddr)

	return ca
}
