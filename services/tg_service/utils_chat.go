package tg_service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func getMessageChatTitle(update *tgbotapi.Update) string {
	if update.Message.Chat != nil {
		if update.Message.Chat.Title != "" {
			return update.Message.Chat.Title
		}

		if update.Message.Chat.UserName != "" {
			return update.Message.Chat.UserName
		}
	}

	return "unknown"
}
