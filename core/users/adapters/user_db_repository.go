package adapters

import (
	"context"
	"errors"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/core/users/domain"
	"github.com/hinha/coai/internal/store/gorm/mysql"
	"gorm.io/gorm"
	"strings"
)

type UserMysqlRepository struct {
	db *mysql.DB
}

func NewUserMysqlRepository(db *mysql.DB) *UserMysqlRepository {
	return &UserMysqlRepository{db: db}
}

func (u *UserMysqlRepository) userTable() *user {
	return Use(u.db.Gorm).User.Table(domain.TableNameUser)
}

func (u *UserMysqlRepository) AddUser(ctx context.Context, user *domain.User) error {
	err := u.userTable().WithContext(ctx).Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserMysqlRepository) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	table := u.userTable()
	user, err := table.WithContext(ctx).Where(table.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}
	// some logic
	return user, nil
}

func (u *UserMysqlRepository) AllUsers(ctx context.Context) ([]query.User, error) {
	table := u.userTable()
	users, err := table.WithContext(ctx).Order(table.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return u.usersModelsToQuery(users)
}

func (u *UserMysqlRepository) usersModelsToQuery(iter []*domain.User) ([]query.User, error) {
	var users []query.User

	for _, user := range iter {
		queryUser := query.User{
			Id:        user.ID,
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		users = append(users, queryUser)
	}

	return users, nil
}

func wrapError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return mysql.RecordNotFound
	case strings.Contains(err.Error(), "1062: Duplicate"):
		return mysql.DataAlreadyExists
	default:
		return err
	}
}
