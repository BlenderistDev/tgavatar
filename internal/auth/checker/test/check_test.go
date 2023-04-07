package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	auth2 "github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"tgavatar/internal/auth/checker"
	"tgavatar/internal/auth/checker/mock_check"
)

func TestCheckerAuthStatus_CheckAuth(t *testing.T) {
	resErr := fmt.Errorf("some err")
	tests := []struct {
		status      *auth2.Status
		expected    bool
		err         error
		expectedErr error
	}{
		{
			status: &auth2.Status{
				Authorized: false,
			},
			expected:    false,
			err:         nil,
			expectedErr: nil,
		},
		{
			status: &auth2.Status{
				Authorized: true,
			},
			expected:    true,
			err:         nil,
			expectedErr: nil,
		},
		{
			status:      nil,
			expected:    false,
			err:         resErr,
			expectedErr: errors.Wrap(resErr, "failed to get auth status for check auth"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, test := range tests {
		ctx := context.Background()

		auth := mock_check.NewMockTgAuthInterface(ctrl)
		auth.EXPECT().Status(ctx).Return(test.status, test.err)

		checker := checker.NewCheckerStatusAuth()
		authorized, err := checker.CheckAuth(ctx, auth)
		assert.Equal(t, test.expected, authorized)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestCheckerAuth_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("some err")
	tests := []struct {
		authorized  bool
		err         error
		expected    bool
		expectedErr error
	}{
		{
			authorized:  true,
			expected:    true,
			err:         nil,
			expectedErr: nil,
		},
		{
			authorized:  false,
			expected:    false,
			err:         nil,
			expectedErr: nil,
		},
		{
			authorized:  false,
			expected:    false,
			err:         resErr,
			expectedErr: errors.Wrap(resErr, "failed to check auth from telegram auth"),
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		tgAuth := &auth2.Client{}

		checkerAuthStatus := mock_check.NewMockCheckerAuthStatusInterface(ctrl)
		checkerAuthStatus.EXPECT().CheckAuth(ctx, tgAuth).Return(test.authorized, test.err)

		checkerAuth := checker.NewCheckerAuth(checkerAuthStatus)
		client := mock_check.NewMockClient(ctrl)
		client.EXPECT().Auth().Return(tgAuth)

		authorized, err := checkerAuth.CheckAuth(ctx, client)

		assert.Equal(t, test.expected, authorized)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}

}

func TestChecker_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("some err")
	tests := []struct {
		err         error
		expected    bool
		expectedErr error
	}{
		{
			expected:    true,
			err:         nil,
			expectedErr: nil,
		},
		{
			expected:    false,
			err:         checker.NoAuthorizedErr,
			expectedErr: nil,
		},
		{
			expected:    false,
			err:         resErr,
			expectedErr: errors.Wrap(resErr, "failed to start TGClient for check auth"),
		},
	}

	for _, test := range tests {

		ctx := context.Background()

		checkFunc := func(ctx context.Context) error {
			return nil
		}

		client := mock_check.NewMockClient(ctrl)
		client.EXPECT().Run(ctx, gomock.AssignableToTypeOf(checkFunc)).Do(func(_ interface{}, f func(ctx context.Context) error) {
			if reflect.ValueOf(checkFunc) != reflect.ValueOf(f) {
				t.Errorf("function mismatch")
			}
		}).Return(test.err)

		telegramFactory := mock_check.NewMockTgFactoryInterface(ctrl)
		telegramFactory.EXPECT().GetClient().Return(client, nil)

		checkerAuth := mock_check.NewMockCheckerAuth(ctrl)
		checkerAuth.EXPECT().GetCheckerFunc(client).Return(checkFunc)

		checker := checker.NewChecker(telegramFactory, checkerAuth)

		res, err := checker.CheckAuth(ctx)
		assert.Equal(t, test.expected, res)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestChecker_CheckAuth_TelegramFactoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("some err")

	ctx := context.Background()

	telegramFactory := mock_check.NewMockTgFactoryInterface(ctrl)
	telegramFactory.EXPECT().GetClient().Return(nil, resErr)

	checkerAuth := mock_check.NewMockCheckerAuth(ctrl)

	checker := checker.NewChecker(telegramFactory, checkerAuth)

	res, err := checker.CheckAuth(ctx)
	assert.False(t, res)
	assert.Equal(t, "failed to create TGClient for check auth: some err", err.Error())

}
