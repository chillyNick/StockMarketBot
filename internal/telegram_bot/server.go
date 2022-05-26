package telegram_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/logger"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/handlers"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/repository"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"google.golang.org/grpc"
)

type Server struct {
	bot        *tgbotapi.BotAPI
	repo       repository.Repository
	grpcClient pb.StockMarketServiceClient
}

func New(tgToken string, repo repository.Repository, conn grpc.ClientConnInterface, debug bool) (*Server, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	return &Server{
		bot:        bot,
		repo:       repo,
		grpcClient: pb.NewStockMarketServiceClient(conn),
	}, nil
}

func (s *Server) Serve() {
	logger.Info.Println("Ready to handle tg messages")

	updateConfig := tgbotapi.NewUpdate(0)
	updates := s.bot.GetUpdatesChan(updateConfig)
	handler := handlers.New(s.repo, s.grpcClient)

	for update := range updates {
		msgConf := handler.HandleUpdate(update)
		if msgConf != nil {
			s.send(msgConf)
		}
	}
}

func (s *Server) send(c tgbotapi.Chattable) (ok bool) {
	_, err := s.bot.Send(c)
	if err != nil {
		logger.Error.Printf("Failed to send msg to telegram: ", err)
		return false
	}

	return true
}
