package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/repository"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type stockMarketServiceServer struct {
	repo repository.Repository
	pb.UnimplementedStockMarketServiceServer
}

func New(repo repository.Repository) *stockMarketServiceServer {
	return &stockMarketServiceServer{repo: repo}
}

var (
	internalError = status.Error(codes.Internal, "Internal error")
)

func (s *stockMarketServiceServer) CreateUser(ctx context.Context, msg *empty.Empty) (*pb.UserId, error) {
	id, err := s.repo.CreateUser(ctx)
	if err != nil {
		logger.Error.Printf("Failed to create user %v\n", err)

		return nil, internalError
	}

	return &pb.UserId{Id: id}, nil
}

func (s *stockMarketServiceServer) GetStocks(ctx context.Context, userId *pb.UserId) (*pb.GetStocksResponse, error) {
	stocks, err := s.repo.GetStocks(ctx, userId.GetId())
	if err != nil {
		logger.Error.Printf("Failed to get stocks %v\n", err)

		return nil, internalError
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
	q, err := getQuote(req.GetStock().GetName())
	if err != nil {
		return nil, err
	}

	stock, err := s.repo.GetStock(ctx, req.GetUserId().GetId(), q.Symbol)
	if errors.Is(err, db.ErrNotFound) {
		err = s.repo.AddStock(ctx, req.GetUserId().GetId(), q.Symbol, req.GetStock().GetAmount(), q.Ask)
		if err != nil {
			logger.Error.Printf("Failed to add stock: %s", err)

			return nil, internalError
		}
		return new(empty.Empty), nil
	}

	if err != nil {
		logger.Error.Printf("Failed to get stock: %s", err)

		return nil, internalError
	}

	amount := stock.Amount + req.GetStock().GetAmount()
	price := (stock.Price*float64(stock.Amount) + q.Ask*float64(req.GetStock().GetAmount())) / float64(amount)
	err = s.repo.UpdateStock(ctx, stock.Id, amount, price)
	if err != nil {
		logger.Error.Printf("Failed to update stock: %s", err)

		return nil, internalError
	}

	return new(empty.Empty), nil
}

func (s *stockMarketServiceServer) RemoveStock(ctx context.Context, req *pb.StockRequest) (*empty.Empty, error) {
	q, err := getQuote(req.GetStock().GetName())
	if err != nil {
		return nil, err
	}

	stock, err := s.repo.GetStock(ctx, req.GetUserId().GetId(), q.Symbol)
	if errors.Is(err, db.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "You don't have stock in the portfolio")
	}

	if err != nil {
		logger.Error.Printf("Failed to get stock: %s", err)

		return nil, internalError
	}

	var action string
	if stock.Amount <= req.GetStock().GetAmount() {
		err = s.repo.RemoveStock(ctx, stock.Id)
		action = "remove"
	} else {
		err = s.repo.UpdateStock(ctx, stock.Id, stock.Amount-req.GetStock().GetAmount(), stock.Price)
		action = "remove"
	}

	if err != nil {
		logger.Error.Printf("Failed to %s stock: %s", action, err)

		return nil, internalError
	}

	return new(empty.Empty), nil
}

func (s *stockMarketServiceServer) AddNotification(ctx context.Context, req *pb.AddNotificationRequest) (*empty.Empty, error) {
	q, err := getQuote(req.GetStockName())
	if err != nil {
		return nil, err
	}

	var nType string
	if q.Bid > req.GetThreshold() {
		nType = models.NotificationTypeDown
	} else {
		nType = models.NotificationTypeUp
	}

	err = s.repo.AddNotification(ctx, req.GetUserId().GetId(), q.Symbol, req.GetThreshold(), nType)
	if err != nil {
		logger.Error.Printf("Failed to add notification: %s", err)

		return nil, internalError
	}

	return new(empty.Empty), nil
}

func (s *stockMarketServiceServer) GetPortfolioChanges(ctx context.Context, req *pb.GetPortfolioChangesRequest) (*pb.GetPortfolioChangesResponse, error) {
	stocks, err := s.repo.GetStocks(ctx, req.GetUserId().GetId())
	if err != nil {
		logger.Error.Printf("Failed to get stocks: %s", err)

		return nil, internalError
	}

	stockByName := make(map[string]models.Stock, len(stocks))
	names := make([]string, 0, len(stocks))
	for _, stock := range stocks {
		stockByName[stock.Name] = stock
		names = append(names, stock.Name)
	}

	res := pb.GetPortfolioChangesResponse{
		Stocks: make([]*pb.StockChanges, 0, len(stocks)),
	}

	curQuotes := quote.List(names)
	for curQuotes.Next() {
		q := curQuotes.Quote()
		stock := stockByName[q.Symbol]

		if req.GetPeriod().String() == "ALL" || stock.CreatedAt.After(getTimeBeforeNow(req.GetPeriod())) {
			res.Stocks = append(res.Stocks, &pb.StockChanges{
				Stock:        &pb.Stock{Name: stock.Name, Amount: stock.Amount},
				OldPrice:     stock.Price,
				CurrentPrice: q.Bid,
			})

			continue
		}

		oldPrice, err := getHistoricalPrice(q.Symbol, req.GetPeriod())
		if err != nil {
			return nil, err
		}

		res.Stocks = append(res.Stocks, &pb.StockChanges{
			Stock:        &pb.Stock{Name: stock.Name, Amount: stock.Amount},
			OldPrice:     oldPrice,
			CurrentPrice: q.Bid,
		})
	}

	if err != nil {
		logger.Error.Printf("Error while getting quotes %s", curQuotes.Err())

		return nil, internalError
	}

	return &res, nil
}

func getTimeBeforeNow(period pb.Period) time.Time {
	var duration time.Duration
	switch period {
	case pb.Period_HOUR:
		duration = time.Hour
	case pb.Period_DAY:
		duration = time.Hour * 24
	case pb.Period_WEEK:
		duration = time.Hour * 24 * 7
	default:
		logger.Error.Fatalf("incorrect period %s", period)
	}

	return time.Now().Add(-duration)
}

func getHistoricalPrice(symbol string, period pb.Period) (float64, error) {
	var interval datetime.Interval
	switch period {
	case pb.Period_HOUR:
		interval = datetime.OneHour
	case pb.Period_DAY:
		interval = datetime.OneDay
	case pb.Period_WEEK:
		interval = "7d"
	default:
		logger.Error.Fatalf("incorrect period %s", period)
	}

	iter := chart.Get(&chart.Params{Symbol: symbol, Interval: interval})
	var oldPrice float64
	for iter.Next() {
		oldPrice = iter.Bar().Open.InexactFloat64()
	}

	if err := iter.Err(); err != nil {
		logger.Error.Printf("Couldn't find old price for stock %s: %s", symbol, err)

		return 0, status.Errorf(codes.NotFound, "Couldn't find old price for stock %s", symbol)
	}

	return oldPrice, nil
}

func getQuote(symbol string) (*finance.Quote, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		logger.Error.Printf("Couldn't get quote by symbol %s: %s", symbol, err)

		return nil, internalError
	}

	if q == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("quote for %s not found", symbol))
	}

	return q, nil
}
