package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "users"

// User mapped from table <user>
type User struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UUID         string         `gorm:"column:uuid;not null" json:"uuid"`
	FirstName    string         `gorm:"column:first_name" json:"first_name"`
	LastName     string         `gorm:"column:last_name" json:"last_name"`
	Email        string         `gorm:"column:email;not null" json:"email"`
	Phone        string         `gorm:"column:phone" json:"phone"`
	Password     string         `gorm:"column:password;not null" json:"password"`
	Intro        string         `gorm:"column:intro" json:"intro"`
	Status       string         `gorm:"column:status;not null;default:inactive" json:"status"`
	Profile      string         `gorm:"column:profile" json:"profile"` // The user details.
	UserGroupsID int64          `gorm:"column:user_groups_id;not null" json:"user_groups_id"`
	LastLogin    time.Time      `gorm:"column:last_login" json:"last_login"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
