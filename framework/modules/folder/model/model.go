package model

import "time"

type Folder struct {
	ID         int        `json:"id,omitempty"`
	UserID     int        `json:"user_id,omitempty"`
	FolderName string     `json:"folder_name,omitempty"`
	Color      string     `json:"color,omitempty"`
	Shared     bool       `json:"shared,omitempty"`
	ParentID   *int       `json:"parent_id,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
