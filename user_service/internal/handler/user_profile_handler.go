package handler

import (
	"context"

	pd "github.com/HunterGooD/voice_friend/contracts/gen/go/user_service"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend/user_service/pkg/server"
)

type UserProfileUsecase interface {
	ChangeAvatar(ctx context.Context)
}

type UserProfileHandler struct {
	pd.UnimplementedUserProfileServer
	uu UserProfileUsecase

	log logger.Logger
}

func NewUserProfileHandler(server *server.GRPCServer, uu UserProfileUsecase, log logger.Logger) {
	userProfileHandler := &UserProfileHandler{uu: uu, log: log}
	pd.RegisterUserProfileServer(server.GetServer(), userProfileHandler)
}

func (up *UserProfileHandler) ChangeAvatar(ctx context.Context, req *pd.AvatarChangeRequest) (*pd.AvatarChangeResponse, error) {
	return nil, nil
}
