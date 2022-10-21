package command_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/core/users/domain/user_groups"
	"github.com/hinha/coai/internal/logger"
	"github.com/hinha/coai/internal/metrics"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCreateUserGroup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name           string
		MockErr        error
		Constructor    func() *domain.UserGroup
		GetConstructor func(group *domain.UserGroup) (*domain.UserGroup, error)

		ShouldFail    bool
		ExpectedError string
	}{
		{
			Name: "create_user_group_status_when_available",
			Constructor: func() *domain.UserGroup {
				return &domain.UserGroup{}
			},
			GetConstructor: func(group *domain.UserGroup) (*domain.UserGroup, error) {
				return group, nil
			},
			MockErr: nil,
		},
		{
			Name: "create_user_group_status_when_not_available",
			Constructor: func() *domain.UserGroup {
				return &domain.UserGroup{}
			},
			MockErr:       errors.New("ERR"),
			ShouldFail:    true,
			ExpectedError: "ERR",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			deps := newDependencies(ctrl)
			construct := tc.Constructor()

			deps.repository.EXPECT().AddUserGroup(gomock.Any(), construct, gomock.Any()).Return(tc.MockErr)

			err := deps.CreateUserGroup.Handle(context.Background(), command.CreateUserGroup{})
			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUpdateUserGroup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		MockErr error

		ShouldFail    bool
		ExpectedError string
	}{
		{
			Name:    "update_user_group",
			MockErr: nil,
		},
		{
			Name:          "update_user_group_status_when_not_found",
			MockErr:       errors.New("ERR"),
			ShouldFail:    true,
			ExpectedError: "ERR",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var anyID int64 = 999

			deps := newDependencies(ctrl)

			deps.repository.EXPECT().UpdateUserGroup(gomock.Any(), &domain.UserGroup{ID: anyID}).Return(tc.MockErr)

			err := deps.UpdateGroup.Handle(context.Background(), command.UpdateUserGroup{
				Id: anyID,
			})
			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDeleteUserGroup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		MockErr error

		ShouldFail    bool
		ExpectedError string
	}{
		{
			Name:    "delete_user_group",
			MockErr: nil,
		},
		{
			Name:          "delete_user_group_status_when_not_exists",
			MockErr:       errors.New("ERR"),
			ShouldFail:    true,
			ExpectedError: "ERR",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var anyID int64 = 999

			deps := newDependencies(ctrl)

			deps.repository.EXPECT().DeleteUserGroup(gomock.Any(), &domain.UserGroup{
				ID: anyID,
				DeletedAt: gorm.DeletedAt(sql.NullTime{
					Valid: true,
				}),
			}).Return(tc.MockErr)

			err := deps.DeleteUserGroup.Handle(context.Background(), command.DeleteUserGroup{
				Id: anyID,
			})
			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestActivateUserGroup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		MockErr error

		ShouldFail    bool
		ExpectedError string
	}{
		{
			Name:    "return_user_group_status_when_active",
			MockErr: nil,
		},
		{
			Name:          "return_user_group_status_when_not_found",
			MockErr:       errors.New("ERR"),
			ShouldFail:    true,
			ExpectedError: "ERR",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var anyID int64 = 999

			deps := newDependencies(ctrl)

			deps.repository.EXPECT().UpdateGroupActive(gomock.Any(), anyID, time.Time{}).Return(tc.MockErr)

			err := deps.ActivateUserGroup.Handle(context.Background(), command.ActivateUserGroup{
				Id:        anyID,
				Timestamp: time.Time{},
			})
			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDeactivateUserGroup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		MockErr error

		ShouldFail    bool
		ExpectedError string
	}{
		{
			Name:    "return_user_group_status_when_inactive",
			MockErr: nil,
		},
		{
			Name:          "return_user_group_status_when_not_found",
			MockErr:       errors.New("ERR"),
			ShouldFail:    true,
			ExpectedError: "ERR",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var anyID int64 = 999

			deps := newDependencies(ctrl)

			deps.repository.EXPECT().UpdateGroupInActive(gomock.Any(), anyID, time.Time{}).Return(tc.MockErr)

			err := deps.DeactivateUserGroup.Handle(context.Background(), command.DeactivateUserGroup{
				Id:        anyID,
				Timestamp: time.Time{},
			})
			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

type dependencies struct {
	repository          *user_groups.MockRepository
	CreateUserGroup     command.CreateUserGroupHandler
	UpdateGroup         command.UpdateUserGroupHandler
	DeleteUserGroup     command.DeleteUserGroupHandler
	ActivateUserGroup   command.ActivateUserGroupHandler
	DeactivateUserGroup command.DeactivateUserGroupHandler
}

func newDependencies(ctrl *gomock.Controller) dependencies {
	repository := user_groups.NewMockRepository(ctrl)
	log := logger.New(logger.Config{})
	metricsClient := metrics.NoOp{}

	return dependencies{
		repository:          repository,
		CreateUserGroup:     command.NewCreateUserGroupHandler(repository, log, metricsClient),
		UpdateGroup:         command.NewUpdateUserGroupHandler(repository, log, metricsClient),
		DeleteUserGroup:     command.NewDeleteUserGroupHandler(repository, log, metricsClient),
		ActivateUserGroup:   command.NewActivateUserGroupHandler(repository, log, metricsClient),
		DeactivateUserGroup: command.NewDeactivateUserGroupHandler(repository, log, metricsClient),
	}
}
