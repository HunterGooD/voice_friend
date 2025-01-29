package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/pkg/errors"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
}

type TokenManager interface {
	GenerateAllTokens(ctx context.Context, uid string, role string) ([]string, error)
	GenerateAccessToken(ctx context.Context, uid string, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid string, role string) (string, error)
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
	var authResponse entity.AuthUserResponse
	if err := uu.ur.AddUser(ctx, user); err != nil {
		uu.log.Error("Add user in db error", map[string]error{
			"error": err,
		})
		return nil, err
	}

	if len(user.UID) == 0 && user.Role == "" {
		uu.log.Error("Invalid params", entity.ErrInternal)
		return nil, errors.Wrap(entity.ErrInternal, "Invalid params")
	}

	tokens, err := uu.tm.GenerateAllTokens(ctx, user.UID.String(), string(user.Role))
	if err != nil {
		uu.log.Error("Error on create jwt tokens", map[string]any{
			"error": err,
		})
		return nil, errors.Wrap(entity.ErrInternal, "error create jwt")
	}

	authResponse = entity.AuthUserResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	}

	return &authResponse, nil
}
