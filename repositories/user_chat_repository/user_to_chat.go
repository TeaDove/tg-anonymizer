package user_chat_repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (r *Repository) ActivateUserToChat(ctx context.Context, userId int64, chatId int64) error {
	err := r.db.
		WithContext(ctx).
		Model(&UserToChat{}).
		Where("user_id = ? AND chat_id = ?", userId, chatId).
		Updates(map[string]any{"status": UserToChatStatusActive, "delete_at": time.Now().UTC().Add(activeChatTtl)}).
		Error
	if err != nil {
		return errors.Wrap(err, "failed to update user to chat")
	}

	_, err = r.GetActiveUserToChat(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "failed to get active user to chat")
	}

	zerolog.Ctx(ctx).Info().
		Int64("user_id", userId).
		Int64("chat_id", chatId).
		Msg("user.to.chat.put")

	return nil
}

func (r *Repository) GetActiveUserToChat(ctx context.Context, userId int64) (int64, error) {
	var userToChat UserToChat

	err := r.db.
		WithContext(ctx).
		First(&userToChat, "user_id = ? AND status = ?", userId, UserToChatStatusActive).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.WithStack(KeyNotFoundErr)
		}

		return 0, errors.Wrap(err, "failed to get active user to chat")
	}

	return userToChat.ChatId, nil
}

func (r *Repository) DelActiveUserToChat(ctx context.Context, userId int64) error {
	err := r.db.
		Model(&UserToChat{}).
		Where("user_id = ? AND status = ?", userId, UserToChatStatusActive).
		Update("status", UserToChatStatusPresent).
		Error
	if err != nil {
		return errors.Wrap(err, "failed to delete user to chat")
	}

	zerolog.Ctx(ctx).Info().
		Int64("userId", userId).
		Msg("user.to.chat.deleted")

	return nil
}
