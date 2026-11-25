package main

import (
	"log/slog"
	"os"

	tmprl "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/tmprl"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	"go.temporal.io/sdk/contrib/sysinfo"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Explicitly enable Debug logs (default SDK level is Warn)
	}))

	opts, err := envconfig.LoadDefaultClientOptions()
	opts.Logger = log.NewStructuredLogger(logger)
	c, err := client.Dial(opts)

	if err != nil {
		logger.Error("Unable to create client", "error", err)
	}
	defer c.Close()

	w := worker.New(c, tmprl.OrderQueue, worker.Options{
		SysInfoProvider: sysinfo.SysInfoProvider(),
	})

	w.RegisterWorkflow(tmprl.PlaceOrderWorkflow)
	w.RegisterActivity(tmprl.GetUserCart)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("Unable to start worker", "error", err)
	}
}
