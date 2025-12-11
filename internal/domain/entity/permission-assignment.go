package entity

import "time"

type PermissionAssignment struct {
	User       User       `json:"user"`
	Permission Permission `json:"permission"`

	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy User      `json:"assigned_by"`

	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy User      `json:"deleted_by"`
}
