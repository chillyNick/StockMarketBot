package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
)

type StockMarketServiceServer struct {
	pb.UnimplementedStockMarketServiceServer
}

func (s *StockMarketServiceServer) FindStock(ctx context.Context, name *pb.StockName) (*pb.Stock, error) {
	fmt.Printf("Find stock method call about %v", name.GetName())

	return &pb.Stock{Name: name}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 6000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStockMarketServiceServer(grpcServer, &StockMarketServiceServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
