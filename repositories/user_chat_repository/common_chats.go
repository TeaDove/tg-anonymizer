package user_chat_repository

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (r *Repository) PutCommonChat(ctx context.Context, userId int64, chat *tgbotapi.Chat) error {
	var alreadyExistedChage UserToChat

	err := r.db.
		WithContext(ctx).
		First(&alreadyExistedChage, "user_id = ? AND chat_id = ?", userId, chat.ID).
		Error
	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "failed to get common chat")
	}

	now := time.Now().UTC()
	err = r.db.
		WithContext(ctx).
		Create(UserToChat{
			UserId:    userId,
			ChatId:    chat.ID,
			ChatTitle: chat.Title,
			Status:    UserToChatStatusPresent,
			CreatedAt: now,
			DeleteAt:  now.Add(presentChatTtl),
		}).
		Error
	if err != nil {
		return errors.Wrap(err, "failed to create user to chat")
	}

	zerolog.Ctx(ctx).Info().Interface("chat", chat).Msg("user.to.chat.added")
	return nil
}

func (r *Repository) GetCommonChats(ctx context.Context, userId int64) ([]UserToChat, error) {
	var userToChats []UserToChat

	err := r.db.
		WithContext(ctx).
		Find(&userToChats, "user_id = ?", userId).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get common chats")
	}

	return userToChats, nil
}
