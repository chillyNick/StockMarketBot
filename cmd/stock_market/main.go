package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/repository/pgx_repository"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/service"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("Error loading .env file %v", err)
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
		logger.Error.Fatal(err)
	}
	defer adp.Close()

	repo := pgx_repository.New(adp)

	go stock_market.TrackNotification(repo, os.Getenv("RABBITMQ_URL"))

	lis, err := net.Listen("tcp", os.Getenv("GRPC_URL"))
	if err != nil {
		logger.Error.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStockMarketServiceServer(grpcServer, service.New(repo))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error.Fatal(err)
	}
}
