package telegram_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/handlers"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"google.golang.org/grpc"
)

type Server struct {
	bot        *tgbotapi.BotAPI
	Repo       Repository
	GrpcClient pb.StockMarketServiceClient
}

func New(tgToken string, repo Repository, conn grpc.ClientConnInterface, debug bool) (*Server, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	return &Server{
		bot:        bot,
		Repo:       repo,
		GrpcClient: pb.NewStockMarketServiceClient(conn),
	}, nil
}

func (s *Server) Serve() {
	logger.Info.Println("Start to handle tg messages")
	updateConfig := tgbotapi.NewUpdate(0)
	updates := s.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		msgConf := handlers.HandleUpdate(s, update)
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
