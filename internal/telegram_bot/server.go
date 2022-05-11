package telegram_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	bot  *tgbotapi.BotAPI
	repo repository
}

func New(tgToken string, repo repository, debug bool) (*server, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	return &server{bot: bot, repo: repo}, nil
}

func (s *server) Serve() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updates := s.bot.GetUpdatesChan(updateConfig)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:6000", opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewStockMarketServiceClient(conn)

	for update := range updates {
		s.handle(update, client)
	}

	return nil
}

func (s *server) send(c tgbotapi.Chattable) {
	if _, err := s.bot.Send(c); err != nil {
		panic(err)
	}
}
