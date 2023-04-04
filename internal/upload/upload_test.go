package upload

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gotd/td/tg"
	mock_upload "tgavatar/internal/upload/mock"
)

func TestUpload_Start(t *testing.T) {
	imgChan := make(chan []byte)
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	loader := mock_upload.NewMockloader(ctrl)

	bytes := make([]byte, 10)
	inputFile := &tg.InputFile{}
	loader.EXPECT().FromBytes(ctx, "avatar.png", bytes).Return(inputFile, nil)

	tgClient := mock_upload.NewMocktgClient(ctrl)

	photosUploadProfilePhotoRequest := &tg.PhotosUploadProfilePhotoRequest{
		File: inputFile,
	}
	const uploadedPhotoId = 123
	uploadRes := &tg.PhotosPhoto{
		Photo: &tg.Photo{ID: uploadedPhotoId},
		Users: nil,
	}
	tgClient.EXPECT().PhotosUploadProfilePhoto(ctx, photosUploadProfilePhotoRequest).Return(uploadRes, nil)

	const photoToDeleteId1 = 456
	const photoToDeleteId2 = 789
	profilePhotos := &tg.PhotosPhotos{
		Photos: []tg.PhotoClass{
			&tg.Photo{ID: photoToDeleteId1},
			&tg.Photo{ID: photoToDeleteId2},
		},
	}
	photosToDelete := []tg.InputPhotoClass{
		&tg.InputPhoto{ID: photoToDeleteId1},
		&tg.InputPhoto{ID: photoToDeleteId2},
	}
	photosGetUserPhotosRequest := &tg.PhotosGetUserPhotosRequest{
		UserID: &tg.InputUserSelf{},
		MaxID:  uploadedPhotoId,
	}
	tgClient.EXPECT().PhotosGetUserPhotos(ctx, photosGetUserPhotosRequest).Return(profilePhotos, nil)
	tgClient.EXPECT().PhotosDeletePhotos(ctx, photosToDelete).Return(nil, nil)

	log := mock_upload.NewMocklog(ctrl)
	log.EXPECT().Info("avatar successfully updated")

	upload := NewUpload(tgClient, loader, log, imgChan)
	go upload.Start(ctx)
	imgChan <- bytes
	time.Sleep(time.Millisecond * 100)
}
