package model

import "time"

type User struct {
	ID           string    `json:"id"`
	NIP          string    `json:"nip"`
	Nama         string    `json:"nama"`
	PasswordHash string    `json:"-"`
	Jabatan      string    `json:"jabatan"`
	Divisi       string    `json:"divisi"`
	Unit         string    `json:"unit"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
