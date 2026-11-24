package main

import (
	"log"

	tmprl "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/tmprl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	"go.temporal.io/sdk/worker"
)

func main() {
	opts, err := envconfig.LoadDefaultClientOptions()
	c, err := client.Dial(opts)

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, tmprl.OrderQueue, worker.Options{})

	w.RegisterWorkflow(tmprl.PlaceOrderWorkflow)
	w.RegisterActivity(tmprl.GetUserCart)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
