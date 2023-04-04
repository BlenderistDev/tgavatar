package upload

import (
	"context"
	"fmt"
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
	time.Sleep(time.Millisecond * 50)
}

func TestUpload_Start_LoadFromBytesError(t *testing.T) {
	resErr := fmt.Errorf("some error")

	imgChan := make(chan []byte)
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	loader := mock_upload.NewMockloader(ctrl)

	bytes := make([]byte, 10)
	loader.EXPECT().FromBytes(ctx, "avatar.png", bytes).Return(nil, resErr)

	tgClient := mock_upload.NewMocktgClient(ctrl)
	log := mock_upload.NewMocklog(ctrl)

	log.EXPECT().Error(gomock.Any()).Do(func(p interface{}) {
		expected := "avatar update error: error while upload file from bytes: some error"
		got := p.(error)
		if got.Error() != expected {
			t.Errorf("expected %v, got: %v", expected, got)
		}
	})

	upload := NewUpload(tgClient, loader, log, imgChan)
	go upload.Start(ctx)
	imgChan <- bytes
	time.Sleep(time.Millisecond * 50)
}

func TestUpload_Start_PhotosUploadProfilePhotoError(t *testing.T) {
	resErr := fmt.Errorf("some error")

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

	tgClient.EXPECT().PhotosUploadProfilePhoto(ctx, photosUploadProfilePhotoRequest).Return(nil, resErr)

	log := mock_upload.NewMocklog(ctrl)
	log.EXPECT().Error(gomock.Any()).Do(func(p interface{}) {
		expected := "avatar update error: error while upload avatar request: some error"
		got := p.(error)
		if got.Error() != expected {
			t.Errorf("expected %v, got: %v", expected, got)
		}
	})

	upload := NewUpload(tgClient, loader, log, imgChan)
	go upload.Start(ctx)
	imgChan <- bytes
	time.Sleep(time.Millisecond * 50)
}

func TestUpload_Start_PhotosGetUserPhotosError(t *testing.T) {
	resErr := fmt.Errorf("some error")

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

	photosGetUserPhotosRequest := &tg.PhotosGetUserPhotosRequest{
		UserID: &tg.InputUserSelf{},
		MaxID:  uploadedPhotoId,
	}
	tgClient.EXPECT().PhotosGetUserPhotos(ctx, photosGetUserPhotosRequest).Return(nil, resErr)

	log := mock_upload.NewMocklog(ctrl)
	log.EXPECT().Error(gomock.Any()).Do(func(p interface{}) {
		expected := "avatar update error: error while deleting old avatars: error while get old avatar request: some error"
		got := p.(error)
		if got.Error() != expected {
			t.Errorf("expected %v, got: %v", expected, got)
		}
	})

	upload := NewUpload(tgClient, loader, log, imgChan)
	go upload.Start(ctx)
	imgChan <- bytes
	time.Sleep(time.Millisecond * 50)
}

func TestUpload_Start_PhotosDeletePhotosError(t *testing.T) {
	resErr := fmt.Errorf("some error")

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
	tgClient.EXPECT().PhotosDeletePhotos(ctx, photosToDelete).Return(nil, resErr)

	log := mock_upload.NewMocklog(ctrl)
	log.EXPECT().Error(gomock.Any()).Do(func(p interface{}) {
		expected := "avatar update error: error while deleting old avatars: error in delete old avatars request: some error"
		got := p.(error)
		if got.Error() != expected {
			t.Errorf("expected %v, got: %v", expected, got)
		}
	})

	upload := NewUpload(tgClient, loader, log, imgChan)
	go upload.Start(ctx)
	imgChan <- bytes
	time.Sleep(time.Millisecond * 50)
}
