package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameRole = "roles"

// Role mapped from table <roles>
type Role struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title        string         `gorm:"column:title;not null" json:"title"`
	Action       string         `gorm:"column:action;not null" json:"action"`
	Description  string         `gorm:"column:description" json:"description"`
	Active       bool           `gorm:"column:active;not null" json:"active"`
	Content      string         `gorm:"column:content" json:"content"`
	MappingRoles []MappingRole  `gorm:"foreignKey:RoleID" json:"mapping_roles"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}
