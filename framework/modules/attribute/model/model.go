package model

import "time"

type Attribute struct {
	ID            int       `json:"id,omitempty"`
	UserID        int       `json:"user_id,omitempty"`
	AttributeName string    `json:"attribute_name,omitempty"`
	AttributeType string    `json:"attribute_type,omitempty"`
	Deleted       bool      `json:"deleted,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	Edges         any       `json:"edges"`
}

type AttributeOption struct {
	ID           int       `json:"id,omitempty"`
	UserID       int       `json:"user_id,omitempty"`
	AttributeID  int       `json:"attribute_id,omitempty"`
	OptionValue  string    `json:"option_value,omitempty"`
	DisplayOrder int       `json:"display_order,omitempty"`
	Deleted      bool      `json:"deleted,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Edges        any       `json:"edges"`
}

type OptionOrder struct {
	OptionID     int `json:"option_id"`
	DisplayOrder int `json:"display_order"`
}
