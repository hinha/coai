package adapters

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"github.com/hinha/coai/internal/store/gorm/mysql/mocks"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"testing"
)

func TestUserMysqlRepository_AddUser(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		wantErr  bool
	}{
		{
			name: "inserted data",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					UUID:         "uuid", // TODO create generator
					FirstName:    "john",
					LastName:     "doe",
					Email:        "john@mail.com",
					Phone:        "0812345678901",
					Password:     "password",
					UserGroupsID: 1,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `users` (`uuid`,`first_name`,`last_name`,`email`,`phone`,`password`,`intro`,`status`,`profile`,`user_groups_id`,`last_login`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 1))

				mock.ExpectCommit()

				return db
			},
		},
		{
			name: "must error inserted data",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					UUID:         "uuid",
					FirstName:    "john",
					LastName:     "doe",
					Email:        "john@mail.com",
					Phone:        "0812345678901",
					Password:     "password",
					UserGroupsID: 1,
				},
			},
			initMock: func() *gorm.DB {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `users` (`uuid`,`first_name`,`last_name`,`email`,`phone`,`password`,`intro`,`status`,`profile`,`user_groups_id`,`last_login`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(int64(0), 1)).WillReturnError(errors.New("ERROR"))

				mock.ExpectCommit()

				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			if err := u.AddUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestUserMysqlRepository_AllUsers(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		initMock func() *gorm.DB
		args     args
		want     []query.User
		wantErr  bool
	}{
		{
			name: "must empty data",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`created_at` DESC")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email"}))

				return db
			},
		},
		{
			name: "must returned data",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`created_at` DESC")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email"}).
						AddRow(1, "uuid", "foo", "doo", "email").
						AddRow(2, "uuid", "foo", "doo", "email"))

				return db
			},
			want: []query.User{{
				Id:        1,
				UUID:      "uuid",
				FirstName: "foo",
				LastName:  "doo",
				Email:     "email",
			}, {
				Id:        2,
				UUID:      "uuid",
				FirstName: "foo",
				LastName:  "doo",
				Email:     "email",
			}},
		},
		{
			name: "should error data not found",
			args: args{
				ctx: context.TODO(),
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`created_at` DESC")).
					WillReturnError(gorm.ErrRecordNotFound)

				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			got, err := u.AllUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserMysqlRepository_GetUser(t *testing.T) {
	db, sql, mock := mocks.NewDatabase()
	defer sql.Close()

	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		args     args
		initMock func() *gorm.DB
		want     *domain.User
		wantErr  bool
	}{
		{
			name: "should returned data",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(1).
					//regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`created_at` DESC")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email"}).
						AddRow(1, "uuid", "foo", "doo", "email"))
				return db
			},
			want: &domain.User{
				ID:        1,
				UUID:      "uuid",
				FirstName: "foo",
				LastName:  "doo",
				Email:     "email",
			},
		},
		{
			name: "should error returned data",
			args: args{
				ctx:    context.TODO(),
				userID: 0,
			},
			initMock: func() *gorm.DB {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email"}).
						AddRow(1, "uuid", "foo", "doo", "email"))
				return db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserMysqlRepository{
				db: &mysql.DB{Gorm: tt.initMock()},
			}
			got, err := u.GetUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserMysqlRepository(t *testing.T) {
	type args struct {
		db *mysql.DB
	}
	tests := []struct {
		name string
		args args
		want *UserMysqlRepository
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &UserMysqlRepository{db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserMysqlRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserMysqlRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
