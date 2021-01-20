package main

import (
	"fmt"
	"github.com/baxiang/soldiers-sortie/lorem/endpoints"
	"github.com/baxiang/soldiers-sortie/lorem/service"
	"github.com/baxiang/soldiers-sortie/lorem/transport"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"

)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	var svc service.Service
	svc = service.LoremService{}
	endpoint := endpoints.Endpoints{
		LoremEndpoint: endpoints.MakeLoremEndpoint(svc),
	}

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		//logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		//logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	r := transport.MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()


	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<- errChan)
}
