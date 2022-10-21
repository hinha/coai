package adapters

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"github.com/hinha/coai/internal/store/gorm/mysql/mocks"
)

func TestNewUserGroupMysqlRepository(t *testing.T) {
	type args struct {
		db *mysql.DB
	}
	tests := []struct {
		name string
		args args
		want *UserGroupMysqlRepository
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &UserGroupMysqlRepository{db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserGroupMysqlRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserGroupMysqlRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGroupMysqlRepository_AddUserGroup(t *testing.T) {
	type args struct {
		ctx      context.Context
		group    *domain.UserGroup
		fnInsert func(*domain.UserGroup) (*domain.UserGroup, error)
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		wantErr  bool
	}{
		{
			name: "must record not found with create data already exists should error",
			args: args{
				ctx:   context.TODO(),
				group: &domain.UserGroup{Name: "ok"},
				fnInsert: func(group *domain.UserGroup) (*domain.UserGroup, error) {
					return nil, mysql.RecordNotFound
				},
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_groups` (`name`,`active`,`updated_at`,`deleted_at`) VALUES (?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "must record not found with create data already exists",
			args: args{
				ctx:   context.TODO(),
				group: &domain.UserGroup{Name: "ok"},
				fnInsert: func(group *domain.UserGroup) (*domain.UserGroup, error) {
					return nil, mysql.RecordNotFound
				},
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_groups` (`name`,`active`,`updated_at`,`deleted_at`) VALUES (?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("1062: Duplicate"))

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `name`=?,`active`=?,`deleted_at`=?,`updated_at`=? WHERE `user_groups`.`name` = ? AND `user_groups`.`deleted_at` IS NOT NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))
				mock.ExpectCommit()
				return db
			},
			wantErr: false,
		},
		{
			name: "existing data with update #2",
			args: args{
				ctx:   context.TODO(),
				group: &domain.UserGroup{Name: "test"},
				fnInsert: func(group *domain.UserGroup) (*domain.UserGroup, error) {
					return nil, nil
				},
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_groups` (`name`,`active`,`updated_at`,`deleted_at`) VALUES (?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("1062: Duplicate"))

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `name`=?,`active`=?,`updated_at`=? WHERE `user_groups`.`name` = ?")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))

				mock.ExpectCommit()
				return db
			},
		},
		{
			name: "existing data with update and should error [DataAlreadyExists] #2",
			args: args{
				ctx:   context.TODO(),
				group: &domain.UserGroup{Name: "test"},
				fnInsert: func(group *domain.UserGroup) (*domain.UserGroup, error) {
					return nil, nil
				},
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_groups` (`name`,`active`,`updated_at`,`deleted_at`) VALUES (?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("1062: Duplicate"))

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `name`=?,`active`=?,`updated_at`=? WHERE `user_groups`.`name` = ?")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "existing data with update and should error not [DataAlreadyExists] #2",
			args: args{
				ctx:   context.TODO(),
				group: &domain.UserGroup{Name: "test"},
				fnInsert: func(group *domain.UserGroup) (*domain.UserGroup, error) {
					return nil, nil
				},
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_groups` (`name`,`active`,`updated_at`,`deleted_at`) VALUES (?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.AddUserGroup(tt.args.ctx, tt.args.group, tt.args.fnInsert); (err != nil) != tt.wantErr {
				t.Errorf("AddUserGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGroupMysqlRepository_GetUserGroup(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx  context.Context
		id   int64
		name string
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		want     *domain.UserGroup
		wantErr  bool
	}{
		{
			name: "must empty data",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `deleted_at` IS NULL AND `user_groups`.`name` = ? ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status"}))

				return db
			},
			wantErr: true,
		},
		{
			name: "must returned and must condition active",
			args: args{
				ctx:  context.TODO(),
				id:   1,
				name: "example",
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `deleted_at` IS NULL AND `user_groups`.`id` = ? AND `user_groups`.`name` = ? ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(1, "example").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active"}).
						AddRow(1, "test", true))

				return db
			},
			want: &domain.UserGroup{ID: 1, Name: "test", Active: true},
		},
		{
			name: "must returned and must id active",
			args: args{
				ctx:  context.TODO(),
				id:   1,
				name: "",
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `deleted_at` IS NULL AND `user_groups`.`id` = ? ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active"}).
						AddRow(1, "test", true))

				return db
			},
			want: &domain.UserGroup{ID: 1, Name: "test", Active: true},
		},
		{
			name: "must returned and must name active",
			args: args{
				ctx:  context.TODO(),
				id:   0,
				name: "example",
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `deleted_at` IS NULL AND `user_groups`.`name` = ? ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active"}).
						AddRow(1, "test", true))

				return db
			},
			want: &domain.UserGroup{ID: 1, Name: "test", Active: true},
		},
		{
			name: "should error any",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `user_groups`.`name` = ? AND `deleted_at` IS NULL ORDER BY `user_groups`.`id` LIMIT 1")).
					WillReturnError(gorm.ErrInvalidTransaction)

				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			got, err := u.GetUserGroup(tt.args.ctx, tt.args.id, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGroupMysqlRepository_UpdateUserGroup(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx   context.Context
		group *domain.UserGroup
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		wantErr  bool
	}{
		{
			name: "should updated data",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 1,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))

				mock.ExpectCommit()
				return db
			},
		},
		{
			name: "should record not found",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 0,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 0))

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "should error internal",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 0,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.UpdateUserGroup(tt.args.ctx, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGroupMysqlRepository_DeleteUserGroup(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx   context.Context
		group *domain.UserGroup
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		wantErr  bool
	}{
		{
			name: "should deleted data",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 1,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `deleted_at`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))

				mock.ExpectCommit()
				return db
			},
		},
		{
			name: "should record not found",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 0,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `deleted_at`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 0))

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "should error internal",
			args: args{
				ctx: context.TODO(),
				group: &domain.UserGroup{
					ID: 0,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `deleted_at`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.DeleteUserGroup(tt.args.ctx, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGroupMysqlRepository_UpdateGroupActive(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx  context.Context
		id   int64
		time time.Time
	}
	tests := []struct {
		name     string
		initMock func() *gorm.DB
		args     args
		wantErr  bool
	}{
		{
			name: "should updated data",
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))

				mock.ExpectCommit()
				return db
			},
		},
		{
			name: "should record not found",
			args: args{
				ctx: context.TODO(),
				id:  0,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 0))

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "should error internal",
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.UpdateGroupActive(tt.args.ctx, tt.args.id, tt.args.time); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGroupActive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGroupMysqlRepository_UpdateGroupInActive(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx  context.Context
		id   int64
		time time.Time
	}
	tests := []struct {
		name     string
		initMock func() *gorm.DB
		args     args
		wantErr  bool
	}{
		{
			name: "should updated data",
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(1), 1))

				mock.ExpectCommit()
				return db
			},
		},
		{
			name: "should record not found",
			args: args{
				ctx: context.TODO(),
				id:  0,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 0))

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
		{
			name: "should error internal",
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_groups` SET `active`=?,`updated_at`=? WHERE `user_groups`.`id` = ? AND `active` = ? AND `user_groups`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidTransaction)

				mock.ExpectCommit()
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.UpdateGroupInActive(tt.args.ctx, tt.args.id, tt.args.time); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGroupInActive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserGroupMysqlRepository_IsActivateGroup(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		want     *domain.UserGroup
		wantErr  bool
	}{
		{
			name: "must empty data",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `user_groups` WHERE `deleted_at` IS NULL AND `user_groups`.`name` = ? ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status"}))

				return db
			},
			wantErr: true,
		},
		{
			name: "must returned",
			args: args{
				ctx:  context.TODO(),
				name: "example",
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT `user_groups`.`active` FROM `user_groups` WHERE `user_groups`.`name` = ? AND `user_groups`.`deleted_at` IS NULL ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs("example").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active"}).
						AddRow(1, "test", true))

				return db
			},
			want: &domain.UserGroup{ID: 1, Name: "test", Active: true},
		},
		{
			name: "must record not found",
			args: args{
				ctx:  context.TODO(),
				name: "example",
			},
			initMock: func() *gorm.DB {
				db, _, mock := mocks.NewDatabase()
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT `user_groups`.`active` FROM `user_groups` WHERE `user_groups`.`name` = ? AND `user_groups`.`deleted_at` IS NULL ORDER BY `user_groups`.`id` LIMIT 1")).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active"}).AddRow(0, "", false).CloseError(gorm.ErrRecordNotFound))
				//WillReturnError(gorm.ErrRecordNotFound)

				return db
			},
			want: &domain.UserGroup{ID: 0, Name: "", Active: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserGroupMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			got, err := u.IsActivateGroup(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsActivateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsActivateGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}
