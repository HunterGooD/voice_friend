package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
)

type UserRepository interface {
	AddUser(ctx context.Context)
}

// TODO: rename functions
type TokenManager interface {
	GenerateAllTokens(ctx context.Context, uid string, roles []string) ([]string, error)
	GenerateAccessToken(ctx context.Context, uid string, roles []string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid string, roles []string) (string, error)
}

type UserUsecase struct {
	ur UserRepository
	tm TokenManager

	log logger.Logger
}

func NewUserUsecase(ur UserRepository, tm TokenManager, log logger.Logger) *UserUsecase {
	return &UserUsecase{ur, tm, log}
}

func (uu *UserUsecase) AddUser(ctx context.Context) {

}
