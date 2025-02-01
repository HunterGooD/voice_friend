package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/pkg/errors"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
	ExistUser(ctx context.Context, login string) (bool, error)
}

type TokenManager interface {
	GenerateAllTokens(ctx context.Context, uid string, role string) ([]string, error)
	GenerateAccessToken(ctx context.Context, uid string, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid string, role string) (string, error)
}

type HashManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) (bool, error)
}

type AuthUsecase struct {
	userRepo UserRepository
	tokenMng TokenManager
	hashMng  HashManager

	log logger.Logger
}

func NewAuthUsecase(ur UserRepository, tm TokenManager, hs HashManager, log logger.Logger) *AuthUsecase {
	return &AuthUsecase{ur, tm, hs, log}
}

func (u *AuthUsecase) RegisterUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error) {

	ok, err := u.userRepo.ExistUser(ctx, user.Login)
	if err != nil {
		u.log.Error("Error checking if user exists", err)
		return nil, errors.Wrap(err, "failed to check user existence")
	}
	if ok {
		return nil, entity.ErrUserAlreadyExists
	}

	// TODO: maybe refactor
	hashPassword, err := u.hashMng.HashPassword(user.Password)
	if err != nil {
		u.log.Error("Error create hash", err)
		return nil, errors.Wrap(entity.ErrInternal, "Error create hash")
	}
	user.Password = hashPassword

	if err := u.userRepo.AddUser(ctx, user); err != nil {
		u.log.Error("Add user in db error", map[string]error{
			"error": err,
		})
		return nil, err
	}

	u.log.Info("Added user in db", map[string]any{
		"login": user.Login,
		"uid":   user.UID.String(),
	})

	tokens, err := u.tokenMng.GenerateAllTokens(ctx, user.UID.String(), string(user.Role))
	if err != nil {
		u.log.Error("Error on create jwt tokens", map[string]any{
			"error": err,
		})
		return nil, errors.Wrap(entity.ErrInternal, "error create jwt")
	}

	return &entity.AuthUserResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	}, nil
}
