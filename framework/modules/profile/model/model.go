package model

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type User struct {
	ID         int        `json:"id,omitempty"`
	Email      string     `json:"email,omitempty"`
	Password   string     `json:"-"`
	Name       string     `json:"name,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	Active     bool       `json:"active,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Avatar     string     `json:"avatar,omitempty"`
	Provider   string     `json:"provider,omitempty"`
	ProviderID string     `json:"provider_id,omitempty"`
	RefCode    *string    `json:"ref_code,omitempty"`
	QrCode     *string    `json:"qr_code,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
}

var (
	ErrPhoneExists = errorString("không thể xử lý yêu cầu")
	ErrEmailExists = errorString("không thể xử lý yêu cầu")
)

type errorString string

func (e errorString) Error() string { return string(e) }

func NewReferenceCode() string {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}

func NewQRCode(ref string) string {
	return "user/" + base64.StdEncoding.EncodeToString([]byte(ref))
}
