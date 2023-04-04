package main

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tgavatar/internal/auth"
	"tgavatar/internal/avatar"
	"tgavatar/internal/cron"
	"tgavatar/internal/log"
	"tgavatar/internal/upload"
	"tgavatar/internal/web"
)

func main() {
	ctx := context.Background()

	logger := log.NewLogger(logrus.New())

	authChecker := auth.Checker{}
	successAuthChan := make(chan struct{})

	authService := auth.NewAuth(ctx, logger, successAuthChan)

	go func() {
		err := web.LaunchAuthServer(authChecker, authService, logger)
		if err != nil {
			panic(errors.Wrap(err, "failed launch auth server"))
		}
	}()

	authorized, err := authChecker.CheckAuth(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed check auth in main"))
	}

	if !authorized {
		logger.Info("wait auth")
		<-successAuthChan
		logger.Info("auth successfully")
	}

	imgChan := make(chan []byte)

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: "storage/session"},
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to create avatar update client"))
	}

	go func() {
		if err := client.Run(ctx, func(ctx context.Context) error {
			loader := uploader.NewUploader(client.API())
			u := upload.NewUpload(client.API(), loader, logger, imgChan)
			u.Start(ctx)

			select {}
		}); err != nil {
			panic(errors.Wrap(err, "failed to run avatar update client"))
		}
	}()

	generator := avatar.NewGenerator(logger)
	err = cron.StartCronAvatarChange(generator, logger, imgChan)
	if err != nil {
		panic(err)
	}

	select {}
}
