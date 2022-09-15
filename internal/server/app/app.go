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
	"github.com/sergalkin/gophkeeper/internal/server/middleware"
	"github.com/sergalkin/gophkeeper/internal/server/service"
	"github.com/sergalkin/gophkeeper/internal/server/storage/postgres"
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

	dbConn, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err = dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	jwtManager, errJwt := jwt.NewJWT(cfg.JWTSecret, cfg.JWTExp)
	if err != nil {
		log.Info(errJwt.Error())
	}

	migrationManager := migrations.NewMigrationManager(cfg)
	err = migrationManager.Up()
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	usersStorage := postgres.NewPostgresUserStorage(dbConn)
	usersGrpcService := service.NewUserGrpc(usersStorage, jwtManager)

	secretTypeStorage := postgres.NewPostgresSecretTypeStorage(dbConn)
	secretTypeGrpcService := service.NewSecretTypeGrpc(secretTypeStorage)

	jwtAuthMiddleware := middleware.NewAuthMiddleware(jwtManager).JwtAuth

	gRPCServer := server.NewGrpcServer(
		server.WithServerConfig(cfg),
		server.WithLogger(log),
		server.WithServices(usersGrpcService, secretTypeGrpcService),
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
