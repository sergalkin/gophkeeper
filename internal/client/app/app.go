package app

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/config"
	"github.com/sergalkin/gophkeeper/internal/client/interceptor"
	"github.com/sergalkin/gophkeeper/internal/client/model"
	"github.com/sergalkin/gophkeeper/internal/client/service"
	"github.com/sergalkin/gophkeeper/internal/client/storage"
	"github.com/sergalkin/gophkeeper/pkg/cert"
	"github.com/sergalkin/gophkeeper/pkg/crypt"
)

type App struct {
	Cancel context.CancelFunc

	SecretService     *service.SecretClientService
	SecretTypeService *service.SecretTypeClientService
	UserService       *service.UserClientService

	Storage storage.Memorier
	Syncer  storage.Syncer
	Cron    *cron.Cron
}

// NewApp - creates Client application.
func NewApp() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	glCtx := model.GlobalContext{Ctx: ctx, Cancel: cancel}

	cfg := config.NewConfig()

	tlsCredential, err := cert.NewSSLConfigService().LoadClientCertificate(cfg)
	if err != nil {
		return nil, fmt.Errorf("error in creating tls creds: %w", err)
	}

	protectedRoutes := map[string]bool{
		"/proto.User/Delete":                   true,
		"/proto.SecretType/GetSecretTypesList": true,
		"/proto.Secret/GetListOfSecretsByType": true,
		"/proto.Secret/CreateSecret":           true,
		"/proto.Secret/GetSecret":              true,
		"/proto.Secret/DeleteSecret":           true,
		"/proto.Secret/Edit":                   true,
	}
	intercept := interceptor.NewAuthInterceptor(protectedRoutes)

	conn, errConn := grpc.Dial(":"+cfg.Port,
		grpc.WithTransportCredentials(tlsCredential),
		grpc.WithUnaryInterceptor(intercept.Unary()),
	)
	if errConn != nil {
		return nil, fmt.Errorf("error in creating grpc con:%w", errConn)
	}

	secretClient := pb.NewSecretClient(conn)
	userClient := pb.NewUserClient(conn)
	secretTypeClient := pb.NewSecretTypeClient(conn)

	cr, errCr := crypt.NewCrypt()
	if errCr != nil {
		return nil, fmt.Errorf("could create crypt")
	}

	memoryStorage := storage.NewMemoryStorage()
	syn := storage.NewSync(memoryStorage, secretClient, &glCtx, cr)

	secretClientService := service.NewSecretClientService(&glCtx, secretClient, memoryStorage, cr, syn)
	userClientService := service.NewUserClientService(&glCtx, userClient)
	secretTypeClientService := service.NewSecretTypeClientService(&glCtx, secretTypeClient)

	c := cron.New()
	c.AddFunc("* * * * *", syn.SyncAll)

	return &App{
		SecretService:     secretClientService,
		SecretTypeService: secretTypeClientService,
		UserService:       userClientService,
		Storage:           memoryStorage,
		Syncer:            syn,
		Cron:              c,
		Cancel:            cancel,
	}, nil
}
