package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
)

type UserRepository interface {
	AddUser(ctx context.Context)
}

type TokenManager interface {
	GenerateAllTokens(ctx context.Context, uid string, roles []string) ([]string, error)
	GenerateAccessToken(ctx context.Context, uid string, roles []string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid string, roles []string) (string, error)
}

type AuthUsecase struct {
	ur UserRepository
	tm TokenManager

	log logger.Logger
}

func NewAuthUsecase(ur UserRepository, tm TokenManager, log logger.Logger) *AuthUsecase {
	return &AuthUsecase{ur, tm, log}
}

func (uu *AuthUsecase) RegisterUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error) {

	return nil, nil
}
