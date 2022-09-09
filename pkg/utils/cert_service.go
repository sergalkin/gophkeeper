package utils

import (
	"crypto/tls"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	"github.com/sergalkin/gophkeeper/internal/client"
	"github.com/sergalkin/gophkeeper/internal/server"
	"github.com/sergalkin/gophkeeper/pkg"
)

var _ SSLConfigLoaderService = (*sslConfigService)(nil)

// SSLConfigLoaderService provides methods for ssl
type SSLConfigLoaderService interface {
	LoadClientCertificate(cfg client.Config) (credentials.TransportCredentials, error)
	LoadServerCertificate(cfg server.Config) (*tls.Config, error)
}

type sslConfigService struct {
	l *zap.Logger
}

// NewSSLConfigService - creates new ssl config service with ability to load client or server certificates.
func NewSSLConfigService() SSLConfigLoaderService {
	return &sslConfigService{
		l: pkg.NewLogger(),
	}
}

// LoadClientCertificate returns client credential TLS by path from client config.
func (s sslConfigService) LoadClientCertificate(cfg client.Config) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(cfg.SSLCertPath, cfg.SSLKeyPath)
	if err != nil {
		s.l.Error(err.Error())
		return nil, err
	}

	return credentials.NewTLS(
		&tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	), nil
}

// LoadServerCertificate returns server tls config by path from server config.
func (s sslConfigService) LoadServerCertificate(cfg server.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.SSLCertPath, cfg.SSLKeyPath)
	if err != nil {
		s.l.Error(err.Error())
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
