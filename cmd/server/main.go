package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/sergalkin/gophkeeper/internal/server/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error in creating app: %w", err))

		return
	}

	application.Server.Start(cancel)

	<-ctx.Done()

	application.Server.Stop()

	application.Logger.Info("Application stopped")
}
