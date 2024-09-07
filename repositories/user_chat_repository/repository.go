package user_chat_repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/pkg/errors"
)

type Repository struct {
	db *gorm.DB
}

var KeyNotFoundErr = errors.New("key not found")

func NewRepository(ctx context.Context, db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}
