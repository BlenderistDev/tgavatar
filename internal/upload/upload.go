package upload

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

// Upload struct for telegram avatar updating
type Upload struct {
	client  *telegram.Client
	imgChan chan []byte
}

// NewUpload constructor for Upload struct
func NewUpload(client *telegram.Client, imgChan chan []byte) Upload {
	return Upload{
		client:  client,
		imgChan: imgChan,
	}
}

// Start run uploading goroutine
func (u Upload) Start(ctx context.Context) {
	for {
		img := <-u.imgChan
		err := u.upload(ctx, img)
		if err != nil {
			log.Println(errors.Wrap(err, "avatar update error"))
		}
		log.Println("avatar successfully updated")
	}
}

func (u Upload) upload(ctx context.Context, img []byte) error {
	loader := uploader.NewUploader(u.client.API())

	file, err := loader.FromBytes(ctx, "avatar.png", img)
	if err != nil {
		return errors.Wrap(err, "error while upload file from bytes")
	}

	res, err := u.client.API().PhotosUploadProfilePhoto(ctx, &tg.PhotosUploadProfilePhotoRequest{
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
	photos, err := u.client.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: &tg.InputUserSelf{},
		Offset: 0,
		MaxID:  maxID,
		Limit:  0,
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

	_, err = u.client.API().PhotosDeletePhotos(ctx, photosToDelete)
	if err != nil {
		return errors.Wrap(err, "error in delete old avatars request")
	}

	return nil
}
