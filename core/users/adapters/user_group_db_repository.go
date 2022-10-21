package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/internal/store/gorm/mysql"
)

type UserGroupMysqlRepository struct {
	db *mysql.DB
}

func NewUserGroupMysqlRepository(db *mysql.DB) *UserGroupMysqlRepository {
	return &UserGroupMysqlRepository{db: db}
}

func (u *UserGroupMysqlRepository) userGroupTable() *userGroup {
	return Use(u.db.Gorm).UserGroup.Table(domain.TableNameUserGroup)
}

func (u *UserGroupMysqlRepository) AddUserGroup(ctx context.Context, group *domain.UserGroup, fnInsert func(*domain.UserGroup) (*domain.UserGroup, error)) error {
	table := u.userGroupTable()
	_, err := fnInsert(group)
	switch wrapError(err) {
	case mysql.RecordNotFound:
		err := wrapError(table.WithContext(ctx).Create(group))
		if errors.Is(err, mysql.DataAlreadyExists) {
			_, err := u.userGroupTable().WithContext(ctx).Unscoped().
				Omit(table.UpdatedAt). // issue by time region availability
				Where(table.Name.Eq(group.Name), table.DeletedAt.IsNotNull()).
				UpdateSimple(
					table.Name.Value(group.Name),
					table.Active.Value(group.Active),
					table.DeletedAt.Null(),
					table.UpdatedAt.Value(group.UpdatedAt),
				)
			return wrapError(err)
		}
		return err
	default:
		err := wrapError(table.WithContext(ctx).Create(group))
		if errors.Is(err, mysql.DataAlreadyExists) {
			_, err := table.WithContext(ctx).
				Omit(table.UpdatedAt). // issue by time region availability
				Where(table.Name.Eq(group.Name)).
				UpdateSimple(
					table.Name.Value(group.Name),
					table.Active.Value(group.Active),
					table.UpdatedAt.Value(group.UpdatedAt),
				)
			return err
		}
		return err
	}
}

func (u *UserGroupMysqlRepository) GetUserGroup(ctx context.Context, id int64, name string) (*domain.UserGroup, error) {
	table := u.userGroupTable()

	db := table.WithContext(ctx).Debug().Unscoped()
	db.Where(table.DeletedAt.Null())
	if id != 0 && name != "" {
		db.Where(table.ID.Eq(id), table.Name.Eq(name))
	} else if id != 0 {
		db.Where(table.ID.Eq(id))
	} else if name != "" {
		db.Where(table.Name.Eq(name))
	} else {
		db.Where(table.Name.Eq(name))
	}
	group, err := db.First()

	if err != nil {
		return nil, wrapError(err)
	}

	return group, nil
}

func (u *UserGroupMysqlRepository) UpdateUserGroup(ctx context.Context, group *domain.UserGroup) error {
	table := u.userGroupTable()
	result, err := table.WithContext(ctx).
		Unscoped().
		Omit(table.UpdatedAt). // issue by time region availability
		Where(table.ID.Eq(group.ID), table.DeletedAt.IsNull()).
		UpdateSimple(table.Active.Value(group.Active), table.UpdatedAt.Value(group.UpdatedAt))

	if err != nil {
		return wrapError(err)
	}

	if result.RowsAffected == 0 {
		return mysql.RecordNotFound
	}
	return nil
}

func (u *UserGroupMysqlRepository) UpdateGroupActive(ctx context.Context, id int64, time time.Time) error {
	table := u.userGroupTable()
	result, err := table.WithContext(ctx).
		Unscoped().
		Omit(table.UpdatedAt).
		Where(table.ID.Eq(id), table.Active.Value(false), table.DeletedAt.IsNull()).
		UpdateSimple(table.Active.Value(true), table.UpdatedAt.Value(time))

	if err != nil {
		return wrapError(err)
	}

	if result.RowsAffected == 0 {
		return mysql.RecordNotFound
	}
	return nil
}

func (u *UserGroupMysqlRepository) UpdateGroupInActive(ctx context.Context, id int64, time time.Time) error {
	table := u.userGroupTable()
	result, err := table.WithContext(ctx).
		Unscoped().
		Omit(table.UpdatedAt).
		Where(table.ID.Eq(id), table.Active.Value(true), table.DeletedAt.IsNull()).
		UpdateSimple(table.Active.Value(false), table.UpdatedAt.Value(time))

	if err != nil {
		return wrapError(err)
	}

	if result.RowsAffected == 0 {
		return mysql.RecordNotFound
	}
	return nil
}

func (u *UserGroupMysqlRepository) DeleteUserGroup(ctx context.Context, group *domain.UserGroup) error {
	table := u.userGroupTable()
	result, err := table.WithContext(ctx).
		Unscoped().
		Omit(table.UpdatedAt, table.DeletedAt). // issue by time region availability
		Where(table.ID.Eq(group.ID), table.DeletedAt.IsNull()).
		UpdateSimple(table.DeletedAt.Value(group.DeletedAt), table.UpdatedAt.Value(group.UpdatedAt))
	if err != nil {
		return wrapError(err)
	}

	if result.RowsAffected == 0 {
		return mysql.RecordNotFound
	}
	return nil
}

func (u *UserGroupMysqlRepository) IsActivateGroup(ctx context.Context, name string) (result *domain.UserGroup, err error) {
	table := u.userGroupTable()
	result, err = table.WithContext(ctx).
		Select(table.Active).
		Where(table.Name.Eq(name)).
		First()

	err = wrapError(err)
	if errors.Is(err, mysql.RecordNotFound) {
		if result == nil {
			result = new(domain.UserGroup)
		}
		result.Active = false
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
