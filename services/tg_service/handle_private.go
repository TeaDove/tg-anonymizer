package tg_service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"tg-anonymizer/repositories/user_chat_repository"

	"github.com/teadove/teasutils/utils/redact_utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) replyWithAvailableChats(ctx context.Context, update *tgbotapi.Update) error {
	chats, err := r.userToChatRepository.GetCommonChats(ctx, update.Message.From.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get common chats")
	}

	if len(chats) == 0 {
		err = r.reply(
			update,
			"Hi! Add me to <italic>desired</italic> chat and sent any message to this chat!",
		)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	text := strings.Builder{}
	text.WriteString("Available chats to choose:\n\n")
	for _, chat := range chats {
		text.WriteString(fmt.Sprintf("<code>%d</code>: %s\n", chat.ChatId, chat.ChatTitle))
	}
	text.WriteString(
		"\nJust pass me ID of chat (number before title, that looks like <code>-1001178533048</code>)",
	)

	err = r.reply(update, text.String())
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) handlePrivateMessageSetChatId(
	ctx context.Context,
	update *tgbotapi.Update,
) error {
	chatId, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		err = r.replyWithAvailableChats(ctx, update)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	err = r.userToChatRepository.ActivateUserToChat(ctx, update.Message.From.ID, chatId)
	if err != nil {
		if errors.Is(err, user_chat_repository.KeyNotFoundErr) {
			err = r.reply(update, "Failed to find chat: <code>%d</code>", chatId)
			if err != nil {
				return errors.Wrap(err, "failed to send message")
			}
			return nil
		}

		return errors.Wrap(err, "failed to put chatId")
	}

	err = r.reply(update, "Chat chosen: <code>%d</code>", chatId)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}
	zerolog.Ctx(ctx).Info().Int64("chat_id", chatId).Msg("user.chose.chat")

	return nil
}

func (r *Service) handlePrivateMessageCommandReset(
	ctx context.Context,
	update *tgbotapi.Update,
) error {
	err := r.userToChatRepository.DelActiveUserToChat(ctx, update.Message.From.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete chat")
	}

	err = r.reply(update, "Active chat reset")
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	err = r.replyWithAvailableChats(ctx, update)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) handlePrivateMessageCommandStart(
	ctx context.Context,
	update *tgbotapi.Update,
) error {
	err := r.replyWithAvailableChats(ctx, update)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) forward(update *tgbotapi.Update, chatId int64) error {
	var err error

	if update.Message.Sticker != nil {
		msg := tgbotapi.NewSticker(chatId, tgbotapi.FileID(update.Message.Sticker.FileID))

		_, err = r.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	if update.Message.Text != "" {
		msg := tgbotapi.NewMessage(chatId, update.Message.Text)

		_, err = r.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	if update.Message.Video != nil {
		msg := tgbotapi.NewVideo(chatId, tgbotapi.FileID(update.Message.Video.FileID))
		msg.Caption = update.Message.Caption

		_, err = r.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	if update.Message.Photo != nil {
		for _, photo := range update.Message.Photo {
			msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(photo.FileID))
			msg.Caption = update.Message.Caption

			_, err = r.bot.Send(msg)
			if err != nil {
				return errors.Wrap(err, "failed to send message")
			}
		}

		return nil
	}

	err = r.reply(update, "I currently only support stickers, photos, videos and regular messages")
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) handlePrivateMessageForward(ctx context.Context, update *tgbotapi.Update) error {
	userId := update.Message.From.ID

	chatId, err := r.userToChatRepository.GetActiveUserToChat(ctx, userId)
	if err != nil {
		if errors.Is(err, user_chat_repository.KeyNotFoundErr) {
			return r.handlePrivateMessageSetChatId(ctx, update)
		}
		return errors.Wrap(err, "failed to get chatId")
	}

	zerolog.Ctx(ctx).
		Debug().
		Int64("chatId", chatId).
		Str("text", redact_utils.RedactWithPrefix(update.Message.Text)).
		Msg("sending.anon.message")

	err = r.forward(update, chatId)
	if err != nil {
		return errors.Wrap(err, "failed to forward")
	}

	return nil
}

func (r *Service) handlePrivateMessageSendMessage(
	ctx context.Context,
	update *tgbotapi.Update,
) error {
	if len(update.Message.Text) <= 5 {
		err := r.reply(update, "Text required")
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}
	}

	err := r.sqsSupplier.SendMessage(ctx, update.Message.Text[5:])
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	err = r.reply(update, "Message sent")
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) handlePrivateMessageReceiveMessage(
	ctx context.Context,
	update *tgbotapi.Update,
) error {
	messages, err := r.sqsSupplier.ReceiveAndDeleteMessage(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to receive message")
	}

	if len(messages) == 0 {
		err = r.reply(update, "No messages in queue")
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

		return nil
	}

	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("Messages: %d\n\n", len(messages)))
	for idx, message := range messages {
		text.WriteString(fmt.Sprintf("%d - %s\n", idx+1, message))
	}

	err = r.reply(update, text.String())
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Service) handlePrivateMessage(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *tgbotapi.Update,
) error {
	defer wg.Done()

	switch update.Message.Command() {
	case "reset":
		return r.handlePrivateMessageCommandReset(ctx, update)
	case "start":
		return r.handlePrivateMessageCommandStart(ctx, update)
	case "send":
		return r.handlePrivateMessageSendMessage(ctx, update)
	case "receive":
		return r.handlePrivateMessageReceiveMessage(ctx, update)
	default:
		return nil
		// return r.handlePrivateMessageForward(ctx, update)
	}
}
