package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client"
	"github.com/sergalkin/gophkeeper/pkg"
	"github.com/sergalkin/gophkeeper/pkg/utils"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cfg := client.NewConfig()
	logger := pkg.NewLogger()

	tlsCredential, err := utils.NewSSLConfigService().LoadClientCertificate(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	conn, errConn := grpc.Dial(":"+cfg.Port, grpc.WithTransportCredentials(tlsCredential))
	if errConn != nil {
		logger.Error(errConn.Error())
	}

	gRPCClient := pb.NewGophKeeperClient(conn)
	go func() {
		fmt.Println(gRPCClient.Hello(ctx, &pb.EmptyRequest{}))
	}()

	releaseResources(ctx, conn, logger)
}

func releaseResources(ctx context.Context, conn *grpc.ClientConn, l *zap.Logger) {
	<-ctx.Done()

	err := conn.Close()
	if err != nil {
		l.Sugar().Errorf("error in closing gRPC client: %s", err.Error())
	} else {
		l.Info("gRPC client was successfully shut down")
	}
}
