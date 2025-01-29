package handler

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend/user_service/pkg/server"
	"github.com/HunterGooD/voice_friend/user_service/pkg/utils"
	pd "github.com/HunterGooD/voice_friend_contracts/gen/go/user_service"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecase interface {
	RegisterUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error)
}

type AuthHandler struct {
	pd.UnimplementedAuthServer
	uu AuthUsecase

	log logger.Logger
}

func NewAuthHandler(server *server.GRPCServer, uu AuthUsecase, log logger.Logger) {
	authHandler := &AuthHandler{uu: uu, log: log}
	pd.RegisterAuthServer(server.GetServer(), authHandler)
}

func (ah *AuthHandler) Register(ctx context.Context, req *pd.RegisterRequest) (*pd.AuthResponse, error) {
	if req.Login == "" || req.Name == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request missing required field: Name or Login or Password")
	}

	if utils.ValidateEmail(*req.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "request email invalid validation")
	}

	if utils.ValidatePhone(*req.Phone) {
		return nil, status.Errorf(codes.InvalidArgument, "request phone invalid validation")
	}

	user := entity.User{
		Login:          req.Login,
		Name:           req.Name,
		Email:          *req.Email,
		Password:       req.Password,
		ProfilePicture: req.ProfilePicture,
		Phone:          req.Phone,
	}
	res, err := ah.uu.RegisterUserUsecase(ctx, &user)
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "request user exists")
		}
		if errors.Is(err, entity.ErrInternal) {
			return nil, status.Errorf(codes.Internal, "request internal error")
		}
	}

	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (ah *AuthHandler) Login(ctx context.Context, req *pd.LoginRequest) (*pd.AuthResponse, error) {
	return nil, nil
}

func (ah *AuthHandler) LogOut(ctx context.Context, req *pd.LogoutRequest) (*pd.LogoutResponse, error) {
	return nil, nil
}
