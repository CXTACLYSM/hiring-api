package entity

import "time"

type PermissionGroup struct {
	Id          string `json:"id"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy User      `json:"created_by"`

	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy User      `json:"updated_by"`

	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy User      `json:"deleted_by"`
}
