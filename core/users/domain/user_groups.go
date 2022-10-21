package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUserGroup = "user_groups"

// UserGroup mapped from table <user_groups>
type UserGroup struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name         string         `gorm:"column:name;not null;index:idx_name,unique" json:"name"`
	Active       bool           `gorm:"column:active;not null" json:"active"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	User         []User         `gorm:"foreignKey:UserGroupsID" json:"users,omitempty"` // Relationship one to many
	MappingRoles []MappingRole  `gorm:"foreignKey:RoleID" json:"mapping_roles"`
}

// TableName UserGroup's table name
func (*UserGroup) TableName() string {
	return TableNameUserGroup
}

func (g *UserGroup) IsAvailable() bool {
	return g.Active == true
}
