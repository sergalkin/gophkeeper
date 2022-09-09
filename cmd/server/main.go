package main

import (
	"context"
	"crypto/tls"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server"
	"github.com/sergalkin/gophkeeper/pkg"
	"github.com/sergalkin/gophkeeper/pkg/utils"
)

type se struct {
	pb.UnimplementedGophKeeperServer
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cfg := server.NewConfig()
	logger := pkg.NewLogger()

	logger.Sugar().Infof("config: %+v", cfg)
	sslConf, err := utils.NewSSLConfigService().LoadServerCertificate(cfg)

	conn, err := tls.Listen("tcp", ":"+cfg.Port, sslConf)
	if err != nil {
		logger.Error(err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGophKeeperServer(grpcServer, se{})

	go func() {
		err = grpcServer.Serve(conn)
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}
	}()

	releaseResources(ctx, grpcServer, logger)
}

func releaseResources(ctx context.Context, grpcServer *grpc.Server, l *zap.Logger) {
	<-ctx.Done()

	grpcServer.Stop()

	l.Info("gRPC server was successfully shut down")
}

func (s se) Hello(ctx context.Context, in *pb.EmptyRequest) (*pb.Response, error) {
	return &pb.Response{Msg: "Privet"}, nil
}
