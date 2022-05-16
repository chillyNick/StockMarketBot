package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/repository"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
}

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

	repo := repository.New(adp)

	go stock_market.TrackNotification(repo, os.Getenv("RABBITMQ_URL"))

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 6000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStockMarketServiceServer(grpcServer, stock_market.New(repo))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
