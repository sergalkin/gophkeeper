package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/sergalkin/gophkeeper/internal/server/config"
	"github.com/sergalkin/gophkeeper/internal/server/service"
	"github.com/sergalkin/gophkeeper/internal/server/storage/postgres"
	"github.com/sergalkin/gophkeeper/pkg/logger"
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

	// TODO добавить миграции
	// TODO добавить логирование запросов
	// TODO посмотреть про интерцепторы? Мидлы?
	// TODO написать логику сохранения пользователя в БД
	// TODO написать логику логина
	// TODO JWT
	// TODO валидацию данных введеных от пользователя
	// TODO шифрование пароля
	usersStorage := postgres.NewPostgresUserStorage(dbConn)
	usersGrpcService := service.NewUserGrpc(usersStorage)

	gRPCServer := server.NewGrpcServer(
		server.WithServerConfig(cfg),
		server.WithLogger(log),
		server.WithServices(usersGrpcService),
	)

	return &App{
		DB:     dbConn,
		Server: gRPCServer,
		Logger: log,
	}, nil
}
