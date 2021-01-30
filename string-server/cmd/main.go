package main

import (
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"

	"github.com/baxiang/soldiers-sortie/string-server/endpoints"
	"github.com/baxiang/soldiers-sortie/string-server/services"
	"github.com/baxiang/soldiers-sortie/string-server/transport"
)

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	svc := &services.StrService{}
	eps := endpoints.MakeEndpoints(svc, logger)
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/palindrome").Handler(transport.GetIsPalHandler(eps.GetIsPalindrome))
	r.Methods(http.MethodGet).Path("/reverse").Handler(transport.GetReverseHandler(eps.GetReverse))
	level.Info(logger).Log("status", "listening", "port", "8080")
	svr := http.Server{
		Addr:    ":8060",
		Handler: r,
	}
	log.Fatal(svr.ListenAndServe())
}
