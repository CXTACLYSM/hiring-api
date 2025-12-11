package entity

import "time"

type PermissionGroupAssignment struct {
	PermissionGroup PermissionGroup `json:"permission_group"`
	User            User            `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy User      `json:"created_by"`

	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy User      `json:"updated_by"`

	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy User      `json:"deleted_by"`
}
