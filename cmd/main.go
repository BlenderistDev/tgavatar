package main

import (
	"context"
	"fmt"

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
	authChecker := auth.Checker{}
	authChan := make(chan struct{})

	authService := auth.NewAuth(authChan)

	go func() {
		err := web.LaunchAuthServer(authChecker, authService)
		if err != nil {
			panic(errors.Wrap(err, "failed launch auth server"))
		}
	}()

	authorized, err := authChecker.CheckAuth(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed check auth in main"))
	}

	if !authorized {
		fmt.Println("wait auth")
		<-authChan
		fmt.Println("auth successfully")
	}

	imgChan := make(chan []byte)

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: "session"},
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to create avatar update client"))
	}

	go func() {
		if err := client.Run(context.Background(), func(ctx context.Context) error {
			u := upload.NewUpload(client, imgChan)
			u.Start()

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
