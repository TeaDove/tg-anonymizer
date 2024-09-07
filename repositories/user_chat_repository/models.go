package user_chat_repository

import "time"

type UserToChat struct {
	UserId    int64
	ChatId    int64
	ChatTitle string
	Status    string
	CreatedAt time.Time
	DeleteAt  time.Time
}

func (UserToChat) TableName() string {
	return "user_to_chat"
}

const (
	UserToChatStatusActive  = "ACTIVE"
	UserToChatStatusPresent = "PRESENT"
)

const (
	presentChatTtl = 24 * time.Hour
	activeChatTtl  = 30 * 24 * time.Hour
)
