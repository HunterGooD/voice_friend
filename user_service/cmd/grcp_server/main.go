package main

import (
	"os"

	"github.com/HunterGooD/voice_friend/user_service/config"
	"github.com/HunterGooD/voice_friend/user_service/internal/auth"
	"github.com/HunterGooD/voice_friend/user_service/internal/handler"
	"github.com/HunterGooD/voice_friend/user_service/internal/repository"
	"github.com/HunterGooD/voice_friend/user_service/internal/usecase"
	"github.com/HunterGooD/voice_friend/user_service/pkg/database"
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend/user_service/pkg/server"
)

func main() {
	// load config
	log := logger.NewJsonLogrusLogger(os.Stdout, os.Getenv("LOG_LEVEL"))
	config_path := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(config_path)
	if err != nil {
		log.Error("Error init config", err)
		panic(err)
	}

	db, err := database.NewPostgresConnection(
		cfg.BuildDSN(),
		cfg.Database.PoolConnection.MaxOpenConns,
		cfg.Database.PoolConnection.MaxIdleConns,
		cfg.Database.PoolConnection.MaxLifeTime,
	)
	if err != nil {
		log.Error("Error init database", err)
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)

	tokenManager, err := auth.NewJWTGenerator(
		cfg.App.CertFilePath,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
		[]string{""},
	)
	if err != nil {
		log.Error("Error init token manager", err)
		panic(err)
	}

	authUsecase := usecase.NewAuthUsecase(userRepository, tokenManager, log)
	userProfileUsecase := usecase.NewUserProfileUsecase(userRepository, tokenManager, log)

	// init gRPC server
	gRPCServer := server.NewGRPCServer(log)

	// register handlers
	handler.NewAuthHandler(gRPCServer, authUsecase, log)
	handler.NewUserProfileHandler(gRPCServer, userProfileUsecase, log)

	if err := gRPCServer.Start(cfg.GetAddress()); err != nil {
		panic(err)
	}
}
