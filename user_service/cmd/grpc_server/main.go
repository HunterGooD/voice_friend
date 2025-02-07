package main

import (
	auth2 "github.com/HunterGooD/voice_friend/user_service/pkg/auth"
	"os"
	"time"

	"github.com/HunterGooD/voice_friend/user_service/config"
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

	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
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

	tokenManager, err := auth2.NewJWTGenerator(
		cfg.App.CertFilePath,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTokenDuration,
		cfg.GetRefreshTokenTime(),
		[]string{""},
	)
	if err != nil {
		log.Error("Error init token manager", err)
		panic(err)
	}

	// TODO: to config params
	hasher := auth2.NewArgon2Hasher(3, 64*1024, 32, 16, 2)

	authUsecase := usecase.NewAuthUsecase(userRepository, tokenManager, hasher)
	userProfileUsecase := usecase.NewUserProfileUsecase(userRepository, tokenManager, log)

	// init gRPC server
	gRPCServer := server.NewGRPCServer(log, 5, time.Duration(30)*time.Second)

	// register handlers
	handler.NewAuthHandler(gRPCServer, authUsecase, log)
	handler.NewUserProfileHandler(gRPCServer, userProfileUsecase, log)

	if err := gRPCServer.Start(cfg.GetAddress()); err != nil {
		panic(err)
	}
}
