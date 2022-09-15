package app

import (
	"context"
	"fmt"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/sergalkin/gophkeeper/internal/server/config"
	"github.com/sergalkin/gophkeeper/internal/server/middleware/auth"
	"github.com/sergalkin/gophkeeper/internal/server/service"
	"github.com/sergalkin/gophkeeper/internal/server/storage/postgres"
	"github.com/sergalkin/gophkeeper/pkg/crypt"
	"github.com/sergalkin/gophkeeper/pkg/jwt"
	"github.com/sergalkin/gophkeeper/pkg/logger"
	"github.com/sergalkin/gophkeeper/pkg/migrations"
	"github.com/sergalkin/gophkeeper/pkg/server"
)

type App struct {
	DB     *pgx.Conn
	Server *server.GrpcServer
	Logger *zap.Logger
}

// NewApp - creates new App.
func NewApp(ctx context.Context) (*App, error) {
	cfg := config.NewConfig()
	log := logger.NewLogger()

	cr, errCr := crypt.NewCrypt()
	if errCr != nil {
		return nil, fmt.Errorf("crypt creating error: %w", errCr)
	}

	dbConn, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err = dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	jwtManager, errJwt := jwt.NewJWT(cfg.JWTSecret, cfg.JWTExp)
	if errJwt != nil {
		return nil, fmt.Errorf("jwtManager creating error: %w", errJwt)
	}

	migrationManager := migrations.NewMigrationManager(cfg)
	err = migrationManager.Up()
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	usersStorage := postgres.NewPostgresUserStorage(dbConn)
	usersGrpcService := service.NewUserGrpc(usersStorage, jwtManager, cr)

	secretTypeStorage := postgres.NewPostgresSecretTypeStorage(dbConn)
	secretTypeGrpcService := service.NewSecretTypeGrpc(secretTypeStorage)

	secretStorage := postgres.NewSecretPostgresStorage(dbConn)
	secretGrpcService := service.NewSecretGrpc(secretStorage)

	//TODO добавление секретов

	jwtAuthMiddleware := auth.NewJwtMiddleware(jwtManager, cr).Auth

	gRPCServer := server.NewGrpcServer(
		server.WithServerConfig(cfg),
		server.WithLogger(log),
		server.WithServices(usersGrpcService, secretTypeGrpcService, secretGrpcService),
		server.WithStreamInterceptors(
			grpczap.StreamServerInterceptor(log),
			grpcauth.StreamServerInterceptor(jwtAuthMiddleware),
			grpcrecovery.StreamServerInterceptor(),
		),
		server.WithUnaryInterceptors(
			grpczap.UnaryServerInterceptor(log),
			grpcauth.UnaryServerInterceptor(jwtAuthMiddleware),
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	return &App{
		DB:     dbConn,
		Server: gRPCServer,
		Logger: log,
	}, nil
}
