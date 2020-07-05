package main

import (
	"net"
	"os"
	"runtime"

	data "github.com/gokusayon/currency/data"
	protos "github.com/gokusayon/currency/protos/currency"
	server "github.com/gokusayon/currency/server"
	"github.com/hashicorp/go-hclog"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := hclog.Default()

	logger.Info("Starting currency server .. ", "os", runtime.GOOS)

	gs := grpc.NewServer()
	er, err := data.NewExchangeRate(logger)

	if err != nil {
		logger.Error("Unable to fetch currency rates", "err", err)
	}

	reflection.Register(gs)

	cs := server.NewCurrency(logger, er)

	protos.RegisterCurrencyServer(gs, cs)

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Error("Unable to listen", "err", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
