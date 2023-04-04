package upload

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=upload.go -destination=./mock/upload.go -package=mock_upload

type log interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

type tgClient interface {
	PhotosUploadProfilePhoto(ctx context.Context, request *tg.PhotosUploadProfilePhotoRequest) (*tg.PhotosPhoto, error)
	PhotosGetUserPhotos(ctx context.Context, request *tg.PhotosGetUserPhotosRequest) (tg.PhotosPhotosClass, error)
	PhotosDeletePhotos(ctx context.Context, id []tg.InputPhotoClass) ([]int64, error)
}

type loader interface {
	FromBytes(ctx context.Context, name string, b []byte) (tg.InputFileClass, error)
}

// Upload struct for telegram avatar updating
type Upload struct {
	client   tgClient
	imgChan  chan []byte
	uploader loader
	log      log
}

// NewUpload constructor for Upload struct
func NewUpload(client tgClient, uploader loader, log log, imgChan chan []byte) Upload {
	return Upload{
		client:   client,
		imgChan:  imgChan,
		uploader: uploader,
		log:      log,
	}
}

// Start run uploading goroutine
func (u Upload) Start(ctx context.Context) {
	for {
		img := <-u.imgChan
		err := u.upload(ctx, img)
		if err != nil {
			u.log.Error(errors.Wrap(err, "avatar update error"))
		}
		u.log.Info("avatar successfully updated")
	}
}

func (u Upload) upload(ctx context.Context, img []byte) error {
	file, err := u.uploader.FromBytes(ctx, "avatar.png", img)
	if err != nil {
		return errors.Wrap(err, "error while upload file from bytes")
	}

	res, err := u.client.PhotosUploadProfilePhoto(ctx, &tg.PhotosUploadProfilePhotoRequest{
		File: file,
	})
	if err != nil {
		return errors.Wrap(err, "error while upload avatar request")
	}

	err = u.deleteOld(ctx, res.GetPhoto().GetID())
	if err != nil {
		return errors.Wrap(err, "error while deleting old avatars")
	}

	return nil
}

func (u Upload) deleteOld(ctx context.Context, maxID int64) error {
	photos, err := u.client.PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: &tg.InputUserSelf{},
		MaxID:  maxID,
	})

	if err != nil {
		return errors.Wrap(err, "error while get old avatar request")
	}

	var photosToDelete []tg.InputPhotoClass
	for _, photo := range photos.GetPhotos() {
		photosToDelete = append(photosToDelete, &tg.InputPhoto{
			ID: photo.GetID(),
		})
	}

	_, err = u.client.PhotosDeletePhotos(ctx, photosToDelete)
	if err != nil {
		return errors.Wrap(err, "error in delete old avatars request")
	}

	return nil
}
