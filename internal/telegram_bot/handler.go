package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

var stocks []string = []string{"apple", "facebook", "amazon", "netflix"}

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	var text string
	switch update.Message.Command() {
	case "start":
		text = "Welcome to the stock market bot. To see all available command type /help"
	case "help":
		text = "/add_stock - Add to the stock portfolio\n/remove_stock - Remove from the stock portfolio"
	case "add_stock":
		text = processAddStockCommand(update)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	sendMessage(bot, msg)
}

func processAddStockCommand(update tgbotapi.Update) string {
	splitMsg := strings.Split(update.Message.Text, " ")
	if len(splitMsg) != 3 {
		return "Type in the next format /add_stock stockName amount"
	}

	stock := ""
	for _, s := range stocks {
		if s == splitMsg[1] {
			stock = splitMsg[1]
		}
	}

	if stock == "" {
		return fmt.Sprintf("Couldn't found stock with name %v", splitMsg[1])
	}

	amount, err := strconv.Atoi(splitMsg[2])
	if err != nil || amount <= 0 {
		return "amount must be a positive number"
	}

	return fmt.Sprintf("%v %v was added into portfolio", amount, splitMsg[1])
}

func sendMessage(bot *tgbotapi.BotAPI, c tgbotapi.Chattable) {
	if _, err := bot.Send(c); err != nil {
		panic(err)
	}
}
