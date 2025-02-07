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
	LoginUserUsecase(ctx context.Context, user *entity.User) (*entity.AuthUserResponse, error)
}

type AuthHandler struct {
	pd.UnimplementedAuthServer
	authUsecase AuthUsecase

	log logger.Logger
}

func NewAuthHandler(gRPCServer *server.GRPCServer, uu AuthUsecase, log logger.Logger) {
	authHandler := &AuthHandler{authUsecase: uu, log: log}
	pd.RegisterAuthServer(gRPCServer.GetServer(), authHandler)
}

func (h *AuthHandler) Register(ctx context.Context, req *pd.RegisterRequest) (*pd.AuthResponse, error) {
	if req.Login == "" || req.Name == "" || req.Password == "" {
		h.log.Warn("Request without params")
		return nil, status.Errorf(codes.InvalidArgument, "Request missing required field: Name or Login or Password")
	}

	if req.Email != nil {
		if utils.ValidateEmail(*req.Email) != true {
			h.log.Warn("Email dont valid", map[string]interface{}{
				"email": *req.Email,
			})
			return nil, status.Errorf(codes.InvalidArgument, "Request email invalid validation")
		}
	}
	if req.Phone != nil {
		if utils.ValidatePhone(*req.Phone) != true {
			h.log.Warn("Phone number dont valid", map[string]interface{}{
				"phone": *req.Phone,
			})
			return nil, status.Errorf(codes.InvalidArgument, "Request phone invalid validation")
		}
	}

	user := entity.User{
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		ProfilePicture: req.ProfilePicture,
		Phone:          req.Phone,
	}
	res, err := h.authUsecase.RegisterUserUsecase(ctx, &user)
	if err != nil {
		//switch errors.Cause(err) {
		//case entity.ErrUserAlreadyExists:
		//	return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		//case entity.ErrInternal:
		//	return nil, status.Errorf(codes.Internal, "internal error")
		//default:
		//	return nil, status.Errorf(codes.Internal, "unknown error %+v", err)
		//}

		if errors.Is(err, entity.ErrUserAlreadyExists) {
			h.log.Warn("User already exists", map[string]interface{}{
				"user":  user,
				"error": err,
			})
			return nil, status.Errorf(codes.AlreadyExists, "Request user exists")
		}
		h.log.Error("Error unknown ", err)
		return nil, status.Errorf(codes.Internal, "Unknown error %+v", err)
	}

	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pd.LoginRequest) (*pd.AuthResponse, error) {
	if req.Login == "" || req.Password == "" {
		h.log.Warn("Request without params")
		return nil, status.Errorf(codes.InvalidArgument, "Request missing required field: Name or Login or Password")
	}

	if req.Email != nil {
		if utils.ValidateEmail(*req.Email) != true {
			h.log.Warn("Email dont valid", map[string]interface{}{
				"email": *req.Email,
			})
			return nil, status.Errorf(codes.InvalidArgument, "Request email invalid validation")
		}
	}
	if req.Phone != nil {
		if utils.ValidatePhone(*req.Phone) != true {
			h.log.Warn("Phone number dont valid", map[string]interface{}{
				"phone": *req.Phone,
			})
			return nil, status.Errorf(codes.InvalidArgument, "Request phone invalid validation")
		}
	}

	user := entity.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}

	h.authUsecase.LoginUserUsecase(ctx, &user)

	return nil, nil
}

func (h *AuthHandler) LogOut(ctx context.Context, req *pd.LogoutRequest) (*pd.LogoutResponse, error) {
	return nil, nil
}
