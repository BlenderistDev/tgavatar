package cron

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"tgavatar/internal/cron/mock_cron"
)

func TestGeneratorJob_generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bytes := []byte("test")
	generator := mock_cron.NewMockgenerator(ctrl)
	generator.EXPECT().Generate(gomock.AssignableToTypeOf(0)).Return(bytes, nil)

	log := mock_cron.NewMocklog(ctrl)

	loc, err := time.LoadLocation("Europe/Moscow")
	assert.Nil(t, err)
	imgChan := make(chan []byte)
	job := generatorJob{
		generator: generator,
		log:       log,
		cron:      nil,
		imgChan:   imgChan,
		loc:       loc,
	}

	var res []byte

	go func() {
		res = <-imgChan
	}()

	job.generate()

	time.Sleep(10)
	assert.Equal(t, res, bytes)

}
