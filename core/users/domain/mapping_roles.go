package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameMappingRole = "mapping_roles"

// MappingRole mapped from table <mapping_roles>
type MappingRole struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RoleID       int64          `gorm:"column:role_id;not null" json:"role_id"`
	UserGroupsID int64          `gorm:"column:user_groups_id;not null" json:"user_groups_id"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName MappingRole's table name
func (*MappingRole) TableName() string {
	return TableNameMappingRole
}
