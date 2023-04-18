package cron

import (
	"fmt"
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

func TestGeneratorJob_generate_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("test")
	generator := mock_cron.NewMockgenerator(ctrl)
	generator.EXPECT().Generate(gomock.AssignableToTypeOf(0)).Return(nil, resErr)

	log := mock_cron.NewMocklog(ctrl)
	log.EXPECT().Error(gomock.Any()).Do(func(err error) {
		assert.Equal(t, "avatar generation failed: test", err.Error())
	})

	loc, err := time.LoadLocation("Europe/Moscow")
	assert.Nil(t, err)
	job := generatorJob{
		generator: generator,
		log:       log,
		cron:      nil,
		imgChan:   nil,
		loc:       loc,
	}

	job.generate()
}
