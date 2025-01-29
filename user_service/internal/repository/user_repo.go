package repository

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewUserRepository(conn *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

func (ur *UserRepository) WithTransaction(tx *sqlx.Tx) *UserRepository {
	return &UserRepository{ur.db, tx}
}

func (ur *UserRepository) AddUser(ctx context.Context, user *entity.User) error {

	return nil
}
