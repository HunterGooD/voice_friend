package handler

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend/user_service/pkg/server"
	pd "github.com/HunterGooD/voice_friend_contracts/gen/go/user_service"
)

type AuthUsecase interface {
	RegisterUserusecase(ctx context.Context)
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
	return nil, nil
}
func (ah *AuthHandler) Login(ctx context.Context, req *pd.LoginRequest) (*pd.AuthResponse, error) {
	return nil, nil
}
func (ah *AuthHandler) LogOut(ctx context.Context, req *pd.LogoutRequest) (*pd.LogoutResponse, error) {
	return nil, nil
}
