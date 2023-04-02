package cron

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"tgavatar/internal/avatar"
)

func StartCronAvatarChange(generator avatar.Generator, imgChan chan []byte) error {
	c := cron.New()
	log.Println("start cron job")
	_, err := c.AddFunc("0 * * * *", func() {
		img, err := generator.Generate(time.Now().Hour())
		if err != nil {
			log.Println(errors.Wrap(err, "avatar generation failed"))
			return
		}
		imgChan <- img
	})

	if err != nil {
		return err
	}

	c.Start()

	return nil
}
