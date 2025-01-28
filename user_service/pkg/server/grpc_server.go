package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	log    logger.Logger
}

func NewGRPCServer(log logger.Logger) *GRPCServer {
	server := grpc.NewServer()
	return &GRPCServer{server, log}
}

func (s *GRPCServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.log.Error("Error init listener: ", map[string]any{
			"error": err,
		})
		return err
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// wait interput signal
		<-quit

		s.log.Info("Initiating graceful shutdown...")
		timer := time.AfterFunc(10*time.Second, func() {
			s.log.Warn("Server couldn't stop gracefully in time. Doing force stop.")
			s.server.Stop()
		})
		defer timer.Stop()

		s.server.GracefulStop()
		s.log.Info("GRPC server stopped")
	}()

	s.log.Info("GRPC server start", map[string]any{
		"grpc_addr": address,
	})
	if err := s.server.Serve(listener); err != nil {
		s.log.Error("Server error: ", map[string]any{
			"error": err,
		})
		return err
	}
	return nil
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
