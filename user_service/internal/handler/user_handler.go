package handler

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
)

type UserUsecase interface {
	AddUser(ctx context.Context)
}

type UserHandler struct {
	uu UserUsecase

	log logger.Logger
}

func NewUserHandler(uu UserUsecase, log logger.Logger) *UserHandler {
	return &UserHandler{uu, log}
}
