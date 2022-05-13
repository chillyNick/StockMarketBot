package main

import (
	"fmt"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/repository"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
)

func main() {
	add := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_STOCK_MARKET_DB"),
	)

	adp, err := db.New(context.Background(), add)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 6000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStockMarketServiceServer(grpcServer, stock_market.New(repository.New(adp)))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
