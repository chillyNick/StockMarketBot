package stock_market

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piquette/finance-go/quote"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type stockMarketServiceServer struct {
	repo repository
	pb.UnimplementedStockMarketServiceServer
}

func New(repo repository) *stockMarketServiceServer {
	return &stockMarketServiceServer{repo: repo}
}

func (s *stockMarketServiceServer) CreateUser(ctx context.Context, msg *empty.Empty) (*pb.UserId, error) {
	id, err := s.repo.CreateUser(ctx)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.UserId{Id: id}, nil
}

func (s *stockMarketServiceServer) GetStocks(ctx context.Context, userId *pb.UserId) (*pb.GetStocksResponse, error) {
	stocks, err := s.repo.GetStocks(ctx, userId.GetId())
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	res := &pb.GetStocksResponse{}
	for _, stock := range stocks {
		res.Stocks = append(res.Stocks, &pb.Stock{
			Name:   stock.Name,
			Amount: stock.Amount,
		})
	}

	return res, nil
}

func (s *stockMarketServiceServer) AddStock(ctx context.Context, req *pb.StockRequest) (*empty.Empty, error) {
	q, err := quote.Get(req.GetName())
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.NotFound, "Incorrect stock name or smth wrong with server")
	}

	stock, err := s.repo.GetStock(ctx, req.GetUserId().GetId(), q.Symbol)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	if errors.Is(err, db.ErrNotFound) {
		err = s.repo.AddStock(ctx, req.GetUserId().GetId(), q.Symbol, req.GetAmount())
		if err != nil {
			log.Println(err)

			return nil, status.Error(codes.Internal, "Internal error")
		}
	} else {
		err = s.repo.UpdateStockAmount(ctx, stock.Id, stock.Amount+req.GetAmount())
	}

	return nil, nil
}

func (s *stockMarketServiceServer) RemoveStock(ctx context.Context, req *pb.StockRequest) (*empty.Empty, error) {
	q, err := quote.Get(req.GetName())
	if err != nil {
		log.Println(err)

		return nil, status.Errorf(codes.NotFound, "Incorrect stock name or smth wrong with server")
	}

	stock, err := s.repo.GetStock(ctx, req.GetUserId().GetId(), q.Symbol)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	if errors.Is(err, db.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "You don't have stock in the portfolio")
	}

	if stock.Amount <= req.GetAmount() {
		err = s.repo.RemoveStock(ctx, stock.Id)
	} else {
		err = s.repo.UpdateStockAmount(ctx, stock.Id, stock.Amount-req.GetAmount())
	}

	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return nil, nil
}
