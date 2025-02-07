package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/pkg/errors"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
	ExistUser(ctx context.Context, login string) (bool, error)
	GetUserPasswordByLogin(ctx context.Context, login string) (string, error)
}

type TokenManager interface {
	GenerateAllTokensAsync(ctx context.Context, uid, role string) ([]string, error)
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
}

func NewAuthUsecase(ur UserRepository, tm TokenManager, hs HashManager) *AuthUsecase {
	return &AuthUsecase{ur, tm, hs}
}

func (u *AuthUsecase) RegisterUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error) {

	ok, err := u.userRepo.ExistUser(ctx, user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "Error check user existence")
	}
	if ok {
		return nil, entity.ErrUserAlreadyExists
	}

	// TODO: maybe refactor
	hashPassword, err := u.hashMng.HashPassword(user.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Error create hash")
	}
	user.Password = hashPassword

	if err := u.userRepo.AddUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "Error create user")
	}

	return u.generateAuthResponse(ctx, user.UID.String(), string(user.Role))
}

func (u *AuthUsecase) LoginUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error) {
	password, err := u.userRepo.GetUserPasswordByLogin(ctx, user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "Error get user password")
	}

	isCorrect, err := u.hashMng.CheckPassword(user.Password, password)
	if err != nil {
		return nil, errors.Wrap(err, "Error check password")
	}

	if !isCorrect {
		return nil, errors.Wrap(entity.ErrInvalidPassword, "Error password not correct")
	}

	return u.generateAuthResponse(ctx, user.UID.String(), string(user.Role))
}

func (u *AuthUsecase) LogoutUserUsecase(ctx context.Context, user *entity.User) error {
	// Deactivate access token and refresh token
	return nil
}

func (u *AuthUsecase) generateAuthResponse(ctx context.Context, uid, role string) (*entity.AuthUserResponse, error) {
	tokens, err := u.tokenMng.GenerateAllTokensAsync(ctx, uid, role)
	if err != nil {
		return nil, errors.Wrap(err, "Error create jwt")
	}

	return &entity.AuthUserResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	}, nil
}
