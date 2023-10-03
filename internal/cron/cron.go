package cron

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

//go:generate mockgen -source=cron.go -destination=./mock_cron/cron.go -package=mock_cron

type log interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

type croner interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Start()
}

type generator interface {
	Generate(hour int) ([]byte, error)
}
type generatorJob struct {
	generator generator
	log       log
	cron      croner
	imgChan   chan []byte
	loc       *time.Location
}

// NewGeneratorJob starts new avatar generator job
func NewGeneratorJob(g generator, l log, i chan []byte, c croner) (*generatorJob, error) {
	job := generatorJob{
		generator: g,
		log:       l,
		cron:      c,
		imgChan:   i,
	}

	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load timezone")
	}

	job.loc = loc

	err = job.startCronAvatarChange()
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// startCronAvatarChange starts cronjob for changing avatar
func (j generatorJob) startCronAvatarChange() error {
	go j.generate()
	j.log.Info("start cron job")
	_, err := j.cron.AddFunc("0 * * * *", j.generate)

	if err != nil {
		return errors.Wrap(err, "failed to start cron job")
	}

	go j.cron.Start()

	return nil
}

func (j generatorJob) generate() {
	img, err := j.generator.Generate(time.Now().In(j.loc).Hour())
	if err != nil {
		j.log.Error(errors.Wrap(err, "avatar generation failed"))
		return
	}
	j.imgChan <- img
}
