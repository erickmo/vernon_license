//go:build wasm

// Package store menyimpan state app di localStorage dan memory.
package store

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// AuthUser adalah user info yang disimpan setelah login berhasil.
type AuthUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

// AuthStore mengelola authentication state menggunakan localStorage.
type AuthStore struct{}

const authKey = "vernon_auth"

// roleWeight memetakan role ke nilai numerik untuk perbandingan hierarki.
var roleWeight = map[string]int{
	"sales":         1,
	"project_owner": 2,
	"superuser":     3,
}

// Save menyimpan auth info ke localStorage.
func (s *AuthStore) Save(user AuthUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	app.Window().Get("localStorage").Call("setItem", authKey, string(data))
	return nil
}

// Load membaca auth info dari localStorage.
// Returns nil jika tidak ada atau data tidak dapat di-parse.
func (s *AuthStore) Load() *AuthUser {
	raw := app.Window().Get("localStorage").Call("getItem", authKey)
	if raw.IsNull() || raw.IsUndefined() {
		return nil
	}
	var user AuthUser
	if err := json.Unmarshal([]byte(raw.String()), &user); err != nil {
		return nil
	}
	if user.Token == "" {
		return nil
	}
	return &user
}

// Clear menghapus auth dari localStorage.
func (s *AuthStore) Clear() {
	app.Window().Get("localStorage").Call("removeItem", authKey)
}

// IsLoggedIn memeriksa apakah ada token valid di localStorage.
func (s *AuthStore) IsLoggedIn() bool {
	u := s.Load()
	return u != nil && u.Token != ""
}

// GetToken mengambil token string. Returns kosong jika tidak login.
func (s *AuthStore) GetToken() string {
	u := s.Load()
	if u == nil {
		return ""
	}
	return u.Token
}

// GetUser mengambil user info. Returns nil jika tidak login.
func (s *AuthStore) GetUser() *AuthUser {
	return s.Load()
}

// GetRole mengambil role user. Returns kosong jika tidak login.
func (s *AuthStore) GetRole() string {
	u := s.Load()
	if u == nil {
		return ""
	}
	return u.Role
}

// HasRole mengecek apakah user punya role yang sama atau lebih tinggi dari minRole.
// Hierarchy: superuser > project_owner > sales.
func (s *AuthStore) HasRole(minRole string) bool {
	role := s.GetRole()
	return roleWeight[role] >= roleWeight[minRole]
}
