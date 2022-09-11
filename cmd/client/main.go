package main

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client"
	"github.com/sergalkin/gophkeeper/pkg/cert"
	"github.com/sergalkin/gophkeeper/pkg/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cfg := client.NewConfig()
	logger := logger.NewLogger()

	tlsCredential, err := cert.NewSSLConfigService().LoadClientCertificate(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	conn, errConn := grpc.Dial(":"+cfg.Port, grpc.WithTransportCredentials(tlsCredential))
	if errConn != nil {
		logger.Error(errConn.Error())
	}

	gRPCClient := pb.NewUserClient(conn)

	go func() {
		logger.Sugar().Info(gRPCClient.Register(ctx, &pb.RegisterRequest{}))
	}()

	releaseResources(ctx, conn, logger)
}

func releaseResources(ctx context.Context, conn *grpc.ClientConn, l *zap.Logger) {
	<-ctx.Done()

	if err := conn.Close(); err != nil {
		l.Sugar().Errorf("error in closing gRPC client: %s", err.Error())

		return
	}

	l.Info("gRPC client was successfully shut down")
}
