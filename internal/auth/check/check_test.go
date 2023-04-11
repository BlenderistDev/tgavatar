package check

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	auth2 "github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"tgavatar/internal/telegram"
)

func TestStatusCheck_CheckAuth(t *testing.T) {
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

		auth := NewMocktgAuthInterface(ctrl)
		auth.EXPECT().Status(ctx).Return(test.status, test.err)

		checker := statusCheck{}
		authorized, err := checker.checkAuth(ctx, auth)
		assert.Equal(t, test.expected, authorized)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestAuthCheck_CheckAuth(t *testing.T) {
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

		statusChecker := NewMockstatusChecker(ctrl)
		statusChecker.EXPECT().checkAuth(ctx, tgAuth).Return(test.authorized, test.err)

		authCheck := authCheck{statusChecker: statusChecker}
		client := NewMockclient(ctrl)
		client.EXPECT().Auth().Return(tgAuth)

		authorized, err := authCheck.checkAuth(ctx, client)

		assert.Equal(t, test.expected, authorized)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}

}

func TestCheck_CheckAuth(t *testing.T) {
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
			err:         noAuthorizedErr,
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

		client := NewMockclient(ctrl)
		client.EXPECT().Run(ctx, gomock.AssignableToTypeOf(checkFunc)).Return(test.err)

		telegramFactory := NewMocktgFactoryInterface(ctrl)
		telegramFactory.EXPECT().GetClient().Return(client, nil)

		authChecker := NewMockauthChecker(ctrl)

		checker := check{
			telegramFactory: telegramFactory,
			authChecker:     authChecker,
		}

		res, err := checker.CheckAuth(ctx)
		assert.Equal(t, test.expected, res)
		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestCheck_CheckAuth_TelegramFactoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("some err")

	ctx := context.Background()

	telegramFactory := NewMocktgFactoryInterface(ctrl)
	telegramFactory.EXPECT().GetClient().Return(nil, resErr)

	authChecker := NewMockauthChecker(ctrl)

	checker := check{
		telegramFactory: telegramFactory,
		authChecker:     authChecker,
	}

	res, err := checker.CheckAuth(ctx)
	assert.False(t, res)
	assert.Equal(t, "failed to create TGClient for check auth: some err", err.Error())

}

func TestTgFactory_GetClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tgClient := &telegram.TGClient{}
	resErr := fmt.Errorf("some err")

	telegramFactory := NewMocktelegramFactory(ctrl)
	telegramFactory.EXPECT().GetClient().Return(tgClient, resErr)

	tgFactory := tgFactory{factory: telegramFactory}

	res, err := tgFactory.GetClient()
	assert.Equal(t, tgClient, res)
	assert.Equal(t, resErr, err)
}

func TestGetCheckerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resErr := fmt.Errorf("some err")
	tests := []struct {
		authorized  bool
		err         error
		expectedErr error
	}{
		{
			authorized:  true,
			err:         nil,
			expectedErr: nil,
		},
		{
			authorized:  false,
			err:         nil,
			expectedErr: noAuthorizedErr,
		},
		{
			authorized:  false,
			err:         resErr,
			expectedErr: errors.Wrap(resErr, "failed to check auth"),
		},
	}

	for _, test := range tests {
		ctx := context.Background()

		client := NewMockclient(ctrl)

		authChecker := NewMockauthChecker(ctrl)
		authChecker.EXPECT().checkAuth(ctx, client).Return(test.authorized, test.err)

		check := check{authChecker: authChecker}
		f := check.getCheckerFunc(client)

		err := f(ctx)

		if test.expectedErr == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestNewChecker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telegramFactory := NewMocktelegramFactory(ctrl)

	assert.IsType(t, check{}, NewChecker(telegramFactory))
}
