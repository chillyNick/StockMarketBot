package stock_market

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
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

		return nil, status.Error(codes.Internal, "Smth wrong with server")
	}

	if q == nil {
		return nil, status.Error(codes.NotFound, "Incorrect stock name")
	}

	stock, err := s.repo.GetStock(ctx, req.GetUserId().GetId(), q.Symbol)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	if errors.Is(err, db.ErrNotFound) {
		err = s.repo.AddStock(ctx, req.GetUserId().GetId(), q.Symbol, req.GetAmount(), q.Ask)
		if err != nil {
			log.Println(err)

			return nil, status.Error(codes.Internal, "Internal error")
		}
	} else {
		amount := stock.Amount + req.GetAmount()
		price := (stock.Price*float64(stock.Amount) + q.Ask*float64(req.GetAmount())) / float64(amount)
		err = s.repo.UpdateStock(ctx, stock.Id, amount, price)
	}

	return new(empty.Empty), nil
}

func (s *stockMarketServiceServer) RemoveStock(ctx context.Context, req *pb.StockRequest) (*empty.Empty, error) {
	q, err := quote.Get(req.GetName())
	if err != nil {
		log.Println(err)

		return new(empty.Empty), status.Errorf(codes.Internal, "Smth wrong with server")
	}

	if q == nil {
		return nil, status.Error(codes.NotFound, "Incorrect stock name")
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
		err = s.repo.UpdateStock(ctx, stock.Id, stock.Amount-req.GetAmount(), stock.Price)
	}

	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return new(empty.Empty), nil
}

var intervalByPeriod = map[pb.Period]datetime.Interval{
	pb.Period_HOUR: datetime.OneHour,
	pb.Period_DAY:  datetime.OneDay,
	pb.Period_WEEK: "7d",
}

var durationByPeriod = map[pb.Period]time.Duration{
	pb.Period_HOUR: time.Hour,
	pb.Period_DAY:  time.Hour * 24,
	pb.Period_WEEK: time.Hour * 24 * 7,
}

func (s *stockMarketServiceServer) GetPortfolioChanges(ctx context.Context, req *pb.GetPortfolioChangesRequest) (*pb.GetPortfolioChangesResponse, error) {
	stocks, err := s.repo.GetStocks(ctx, req.GetUserId().GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
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

		if req.GetPeriod().String() == "ALL" || stock.CreatedAt.After(time.Now().Add(-durationByPeriod[req.GetPeriod()])) {
			res.Stocks = append(res.Stocks, &pb.StockChanges{
				Stock:        &pb.Stock{Name: stock.Name, Amount: stock.Amount},
				OldPrice:     stock.Price,
				CurrentPrice: q.Bid,
			})

			continue
		}

		params := &chart.Params{
			Symbol:   stock.Name,
			Interval: intervalByPeriod[req.GetPeriod()],
		}
		iter := chart.Get(params)
		var oldPrice float64
		for iter.Next() {
			oldPrice = iter.Bar().Open.InexactFloat64()
		}
		if err := iter.Err(); err != nil {
			log.Println(err)

			return nil, status.Errorf(codes.NotFound, "Couldn't find old price for stock: %v", stock.Name)
		}

		res.Stocks = append(res.Stocks, &pb.StockChanges{
			Stock:        &pb.Stock{Name: stock.Name, Amount: stock.Amount},
			OldPrice:     oldPrice,
			CurrentPrice: q.Bid,
		})
	}

	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &res, nil
}

func (s *stockMarketServiceServer) AddNotification(ctx context.Context, req *pb.AddNotificationRequest) (*empty.Empty, error) {
	q, err := quote.Get(req.GetStockName())
	if err != nil {
		log.Println(err)

		return nil, status.Errorf(codes.Internal, "Smth wrong with server")
	}

	if q == nil {
		return nil, status.Error(codes.NotFound, "Incorrect stock name")
	}

	var nType string
	if q.Bid > req.GetThreshold() {
		nType = models.NotificationTypeDown
	} else {
		nType = models.NotificationTypeUp
	}

	err = s.repo.AddNotification(ctx, req.GetUserId().GetId(), req.GetStockName(), req.GetThreshold(), nType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Smth wrong with server")
	}

	return new(empty.Empty), nil
}
