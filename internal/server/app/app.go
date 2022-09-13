package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/sergalkin/gophkeeper/internal/server/config"
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

	// TODO добавить проверку на наличие jwt валидного токена при получении списка типов секретов
	// TODO посмотреть про интерцепторы? Мидлы?
	// TODO валидацию данных введеных от пользователя
	// TODO прокидывание токена между клиентом и сервером через ctx?
	// TODO добавить логирование запросов?
	usersStorage := postgres.NewPostgresUserStorage(dbConn)
	usersGrpcService := service.NewUserGrpc(usersStorage, jwtManager)

	secretTypeStorage := postgres.NewPostgresSecretTypeStorage(dbConn)
	secretTypeGrpcService := service.NewSecretTypeGrpc(secretTypeStorage, jwtManager)

	gRPCServer := server.NewGrpcServer(
		server.WithServerConfig(cfg),
		server.WithLogger(log),
		server.WithServices(usersGrpcService, secretTypeGrpcService),
	)

	return &App{
		DB:     dbConn,
		Server: gRPCServer,
		Logger: log,
	}, nil
}
