package upload

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

type Upload struct {
	client  *telegram.Client
	imgChan chan []byte
}

func NewUpload(client *telegram.Client, imgChan chan []byte) Upload {
	return Upload{
		client:  client,
		imgChan: imgChan,
	}
}

func (u Upload) Start() {
	for {
		img := <-u.imgChan
		err := u.upload(img)
		if err != nil {
			panic(err)
		}
	}
}

func (u Upload) upload(img []byte) error {
	loader := uploader.NewUploader(u.client.API())

	file, err := loader.FromBytes(context.Background(), "avatar.png", img)
	if err != nil {
		return err
	}

	res, err := u.client.API().PhotosUploadProfilePhoto(context.Background(), &tg.PhotosUploadProfilePhotoRequest{
		File: file,
	})
	if err != nil {
		return err
	}

	fmt.Println(res.String())

	err = u.deleteOld(context.Background(), res.GetPhoto().GetID())
	if err != nil {
		return err
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
		return err
	}

	var photosToDelete []tg.InputPhotoClass
	for _, photo := range photos.GetPhotos() {
		photosToDelete = append(photosToDelete, &tg.InputPhoto{
			ID: photo.GetID(),
		})
	}

	_, err = u.client.API().PhotosDeletePhotos(ctx, photosToDelete)
	if err != nil {
		return err
	}

	return nil
}
