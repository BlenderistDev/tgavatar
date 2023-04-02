package main

import (
	"context"
	"log"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"tgavatar/internal/auth"
	"tgavatar/internal/avatar"
	"tgavatar/internal/cron"
	"tgavatar/internal/upload"
	"tgavatar/internal/web"
)

func main() {
	ctx := context.Background()

	authChecker := auth.Checker{}
	successAuthChan := make(chan struct{})

	authService := auth.NewAuth(ctx, successAuthChan)

	go func() {
		err := web.LaunchAuthServer(authChecker, authService)
		if err != nil {
			panic(errors.Wrap(err, "failed launch auth server"))
		}
	}()

	authorized, err := authChecker.CheckAuth(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed check auth in main"))
	}

	if !authorized {
		log.Println("wait auth")
		<-successAuthChan
		log.Println("auth successfully")
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
			u := upload.NewUpload(client, imgChan)
			u.Start(ctx)

			select {}
		}); err != nil {
			panic(errors.Wrap(err, "failed to run avatar update client"))
		}
	}()

	err = cron.StartCronAvatarChange(avatar.Generator{}, imgChan)
	if err != nil {
		panic(err)
	}

	select {}
}
