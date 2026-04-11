package model

import (
	"encoding/xml"
	"time"
)

type LanguageResourceDTO struct {
	ID         int       `json:"id,omitempty" xml:"-"`
	LanguageID int       `json:"language_id,omitempty" xml:"-"`
	Key        string    `json:"key,omitempty" xml:"key,attr"`
	Value      string    `json:"value,omitempty" xml:",chardata"`
	CreatedAt  time.Time `json:"created_at,omitempty" xml:"-"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" xml:"-"`
}

type LanguageDTO struct {
	ID         int                    `json:"id,omitempty" xml:"-"`
	Code       string                 `json:"code,omitempty" xml:"code,attr"`
	Name       string                 `json:"name,omitempty" xml:"name,attr"`
	NativeName string                 `json:"native_name,omitempty" xml:"native_name,attr"`
	IsDefault  bool                   `json:"is_default,omitempty" xml:"is_default,attr"`
	Active     bool                   `json:"active,omitempty" xml:"active,attr"`
	CreatedAt  time.Time              `json:"created_at,omitempty" xml:"-"`
	UpdatedAt  time.Time              `json:"updated_at,omitempty" xml:"-"`
	Resources  []*LanguageResourceDTO `json:"resources,omitempty" xml:"resources>resource,omitempty"`
}

type LanguageXMLDocument struct {
	XMLName    xml.Name               `xml:"language"`
	Code       string                 `xml:"code,attr,omitempty"`
	Name       string                 `xml:"name,attr,omitempty"`
	NativeName string                 `xml:"native_name,attr,omitempty"`
	Resources  []*LanguageResourceDTO `xml:"resources>resource"`
}

type LanguageOptionDTO struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
	IsDefault  bool   `json:"is_default"`
}

type CurrentUserLanguageDTO struct {
	SelectedCode  *string            `json:"selected_code,omitempty"`
	EffectiveCode string             `json:"effective_code"`
	Language      *LanguageOptionDTO `json:"language,omitempty"`
}

type AdminResourcesDTO struct {
	RequestedCode string             `json:"requested_code,omitempty"`
	EffectiveCode string             `json:"effective_code"`
	Language      *LanguageOptionDTO `json:"language,omitempty"`
	Resources     map[string]string  `json:"resources"`
}
