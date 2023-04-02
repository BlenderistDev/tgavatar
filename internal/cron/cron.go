package cron

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"tgavatar/internal/avatar"
)

// StartCronAvatarChange starts cronjob for changing avatar
func StartCronAvatarChange(generator avatar.Generator, imgChan chan []byte) error {
	c := cron.New()
	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		return errors.Wrap(err, "failed to load timezone")
	}
	log.Println("start cron job")
	_, err = c.AddFunc("0 * * * *", func() {
		img, err := generator.Generate(time.Now().In(loc).Hour())
		if err != nil {
			log.Println(errors.Wrap(err, "avatar generation failed"))
			return
		}
		imgChan <- img
	})

	if err != nil {
		return errors.Wrap(err, "failed to start cron job")
	}

	c.Start()

	return nil
}
