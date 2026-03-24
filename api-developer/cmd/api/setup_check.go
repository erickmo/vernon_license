//go:build !wasm

// Package main — setup_check.go menjalankan setup wizard sebelum FX dimulai.
// Jika .env belum ada, server ini serve form wizard yang mengurus segalanya:
// buat DB, jalankan migrasi, buat superuser, tulis .env, lalu restart.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/crypto/bcrypt"

	"github.com/flashlab/vernon-license/migrations"
)

// isSetupRequired returns true jika .env belum ada ATAU koneksi DB gagal.
func isSetupRequired() bool {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return true
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return true
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return true
	}
	defer db.Close()
	db.SetConnMaxLifetime(3 * time.Second)
	if err := db.Ping(); err != nil {
		return true
	}
	return false
}

// serveSetupWizard menjalankan HTTP server minimal yang menyajikan form setup wizard.
// Blocking — tidak return sampai proses di-replace via syscall.Exec.
func serveSetupWizard(port string) {
	if port == "" {
		port = "8081"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /setup/test-db", handleTestDB)
	mux.HandleFunc("POST /setup/apply", handleSetupApply)
	mux.HandleFunc("/", serveWizardHTML)

	fmt.Printf("Vernon License — Setup Wizard berjalan di http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "setup server error: %v\n", err)
		os.Exit(1)
	}
}

// serveWizardHTML menyajikan web/setup.html.
func serveWizardHTML(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("web/setup.html")
	if err != nil {
		http.Error(w, "web/setup.html tidak ditemukan. Pastikan menjalankan dari direktori api-developer/", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content)
}

// setupRequest adalah payload JSON dari form wizard.
type setupRequest struct {
	DBHost        string `json:"db_host"`
	DBPort        string `json:"db_port"`
	DBName        string `json:"db_name"`
	DBUser        string `json:"db_user"`
	DBPassword    string `json:"db_password"`
	AppPort       string `json:"app_port"`
	JWTSecret     string `json:"jwt_secret"`
	AdminName     string `json:"admin_name"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
}

// validate memeriksa semua field wajib.
func (req *setupRequest) validate() error {
	if req.DBHost == "" || req.DBPort == "" || req.DBName == "" || req.DBUser == "" {
		return fmt.Errorf("DB host, port, nama, dan user wajib diisi")
	}
	if len(req.JWTSecret) < 32 {
		return fmt.Errorf("JWT secret minimal 32 karakter")
	}
	if req.AdminName == "" || req.AdminEmail == "" || req.AdminPassword == "" {
		return fmt.Errorf("nama, email, dan password admin wajib diisi")
	}
	if len(req.AdminPassword) < 8 {
		return fmt.Errorf("password admin minimal 8 karakter")
	}
	if req.AppPort == "" {
		req.AppPort = "8081"
	}
	return nil
}

// dbDSN membangun PostgreSQL DSN. Gunakan dbName="" untuk target DB dari request.
func (req *setupRequest) dbDSN(dbName string) string {
	if dbName == "" {
		dbName = req.DBName
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		req.DBUser, req.DBPassword, req.DBHost, req.DBPort, dbName)
}

// handleTestDB mencoba koneksi ke PostgreSQL dan mengembalikan status.
func handleTestDB(w http.ResponseWriter, r *http.Request) {
	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request tidak valid"})
		return
	}

	db, err := sql.Open("postgres", req.dbDSN("postgres"))
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)
	if err := db.Ping(); err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "error", "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "Koneksi berhasil"})
}

// handleSetupApply menjalankan setup lengkap:
// 1. Buat database jika belum ada
// 2. Jalankan migrasi
// 3. Buat superuser
// 4. Tulis .env
// 5. Restart proses
func handleSetupApply(w http.ResponseWriter, r *http.Request) {
	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request tidak valid"})
		return
	}

	if err := req.validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// — Step 1: Konek ke DB admin (postgres), buat target DB jika belum ada —
	adminDB, err := sql.Open("postgres", req.dbDSN("postgres"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal membuka koneksi: " + err.Error()})
		return
	}
	defer adminDB.Close()
	adminDB.SetConnMaxLifetime(10 * time.Second)

	if err := adminDB.Ping(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal konek ke PostgreSQL: " + err.Error()})
		return
	}

	var dbExists bool
	_ = adminDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", req.DBName).Scan(&dbExists)
	if !dbExists {
		// CREATE DATABASE tidak mendukung parameterized query — identifier di-sanitize manual
		safeName := sanitizeIdentifier(req.DBName)
		if _, err := adminDB.Exec("CREATE DATABASE " + safeName); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal membuat database: " + err.Error()})
			return
		}
	}

	// — Step 2: Konek ke target DB, grant schema, jalankan migrasi —
	targetDB, err := sql.Open("postgres", req.dbDSN(""))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal konek ke database: " + err.Error()})
		return
	}
	defer targetDB.Close()

	if err := targetDB.Ping(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal ping database: " + err.Error()})
		return
	}

	// GRANT tidak mendukung parameterized query — identifier di-sanitize manual
	safeUser := sanitizeIdentifier(req.DBUser)
	_, _ = targetDB.Exec("GRANT ALL ON SCHEMA public TO " + safeUser)

	src := migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations.FS,
		Root:       ".",
	}
	if _, err := migrate.Exec(targetDB, "postgres", src, migrate.Up); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "migrasi gagal: " + err.Error()})
		return
	}

	// — Step 3: Buat superuser jika belum ada user sama sekali —
	userNote := "✓ Superuser dibuat"
	var userCount int
	_ = targetDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), 12)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal hash password: " + err.Error()})
			return
		}
		now := time.Now().UTC()
		if _, err = targetDB.Exec(
			`INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, 'superuser', true, $5, $5)`,
			uuid.New().String(), req.AdminName, req.AdminEmail, string(hash), now,
		); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal membuat superuser: " + err.Error()})
			return
		}
	} else {
		userNote = fmt.Sprintf("⚠ User sudah ada (%d), skip pembuatan superuser", userCount)
	}

	// — Step 4: Tulis .env —
	if err := os.WriteFile(".env", []byte(buildEnvContent(req)), 0600); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal menulis .env: " + err.Error()})
		return
	}

	// Return success ke browser sebelum restart
	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"message":   "Setup selesai! Server sedang restart...",
		"port":      req.AppPort,
		"user_note": userNote,
	})

	// Restart proses dengan config baru — .env sudah ada, main() akan load normal
	go func() {
		time.Sleep(300 * time.Millisecond)
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "restart gagal: %v\n", err)
			return
		}
		if err := syscall.Exec(exe, os.Args, os.Environ()); err != nil {
			fmt.Fprintf(os.Stderr, "exec gagal: %v\n", err)
		}
	}()
}

// buildEnvContent membuat isi file .env dari setupRequest.
func buildEnvContent(req setupRequest) string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		req.DBUser, req.DBPassword, req.DBHost, req.DBPort, req.DBName)
	return strings.Join([]string{
		"DATABASE_URL=" + dsn,
		"JWT_SECRET=" + req.JWTSecret,
		"PORT=" + req.AppPort,
		"LOG_LEVEL=info",
		"STORAGE_PATH=./storage",
		"LICENSE_CHECK_INTERVAL=6h",
		"COMPANY_NAME=Vernon",
		"COMPANY_ADDRESS=",
		"COMPANY_PHONE=",
		"COMPANY_EMAIL=",
		"COMPANY_LOGO_PATH=",
		"",
	}, "\n")
}

// sanitizeIdentifier menghapus karakter berbahaya dari SQL identifier (hanya huruf, angka, underscore).
// Digunakan untuk DDL statements yang tidak mendukung parameterized query.
func sanitizeIdentifier(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// writeJSON menulis response JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
