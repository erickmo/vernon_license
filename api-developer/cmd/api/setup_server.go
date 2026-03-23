package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	devmigrations "github.com/flashlab/flasherp-developer-api/migrations"
)

//go:embed web/setup.html
var setupPageHTML string

// ---------- job store untuk SSE progress ----------

type progressEvent struct {
	Step    string `json:"step,omitempty"`
	Percent int    `json:"percent,omitempty"`
	Done    bool   `json:"done,omitempty"`
	OK      bool   `json:"ok,omitempty"`
	Error   string `json:"error,omitempty"`
	APIURL  string `json:"api_url,omitempty"`
	Email   string `json:"email,omitempty"`
}

type installJob struct {
	ch chan progressEvent
}

var (
	jobs   = map[string]*installJob{}
	jobsMu sync.Mutex
)

func newJob() (string, *installJob) {
	id := generateSecret(8)
	job := &installJob{ch: make(chan progressEvent, 32)}
	jobsMu.Lock()
	jobs[id] = job
	jobsMu.Unlock()
	return id, job
}

func getJob(id string) (*installJob, bool) {
	jobsMu.Lock()
	defer jobsMu.Unlock()
	j, ok := jobs[id]
	return j, ok
}

func deleteJob(id string) {
	jobsMu.Lock()
	delete(jobs, id)
	jobsMu.Unlock()
}

// ---------- pre-flight ----------

// needsSetup returns true jika .env tidak ada, DATABASE_URL kosong, atau DB tidak bisa dikoneksi.
func needsSetup() bool {
	dbURL := readEnvValue(".env", "DATABASE_URL")
	if dbURL == "" {
		return true
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return true
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx) != nil
}

// ---------- setup server ----------

func runSetupServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, setupPageHTML)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/setup", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})
	mux.HandleFunc("/setup/api/status", handleSetupStatus)
	mux.HandleFunc("/setup/api/test-connection", handleTestConnection)
	mux.HandleFunc("/setup/api/install", handleInstallStart)
	mux.HandleFunc("/setup/api/install/stream", handleInstallStream)

	fmt.Printf("\n┌──────────────────────────────────────────────┐\n")
	fmt.Printf("│  Vernon License API — Setup Required          │\n")
	fmt.Printf("│                                               │\n")
	fmt.Printf("│  Buka browser dan kunjungi:                   │\n")
	fmt.Printf("│  http://localhost:%s/setup                 │\n", port)
	fmt.Printf("│                                               │\n")
	fmt.Printf("│  Setelah setup selesai, restart server ini.   │\n")
	fmt.Printf("└──────────────────────────────────────────────┘\n\n")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "setup server error: %v\n", err)
	}
}

// ---------- handlers ----------

func handleSetupStatus(w http.ResponseWriter, r *http.Request) {
	dbURL := readEnvValue(".env", "DATABASE_URL")
	if dbURL == "" {
		jsonResp(w, map[string]any{"installed": false})
		return
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		jsonResp(w, map[string]any{"installed": false})
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var count int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE role = 'superuser'`).Scan(&count); err != nil {
		jsonResp(w, map[string]any{"installed": false})
		return
	}
	port := readEnvValue(".env", "HTTP_PORT")
	if port == "" {
		port = "8081"
	}
	jsonResp(w, map[string]any{
		"installed": count > 0,
		"api_url":   "http://localhost:" + port,
	})
}

// handleTestConnection:
// - Koneksi ke PostgreSQL server (melalui database "postgres")
// - Jika database target tidak ada → OK (akan dibuat saat install)
// - Jika database ada tapi sudah punya tabel → error
// - Jika database ada dan kosong → OK
func handleTestConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req dbConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResp(w, map[string]any{"ok": false, "error": "request tidak valid"})
		return
	}
	if req.Name == "" {
		jsonResp(w, map[string]any{"ok": false, "error": "nama database wajib diisi"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	// Cek server PostgreSQL via database "postgres"
	adminDSN := buildDSN(req.Host, req.Port, "postgres", req.User, req.Password)
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		jsonResp(w, map[string]any{"ok": false, "error": "DSN tidak valid: " + err.Error()})
		return
	}
	defer adminDB.Close()

	if err := adminDB.PingContext(ctx); err != nil {
		jsonResp(w, map[string]any{"ok": false, "error": "Tidak dapat terhubung ke server PostgreSQL: " + err.Error()})
		return
	}

	// Cek apakah database target sudah ada
	var dbExists bool
	_ = adminDB.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, req.Name).
		Scan(&dbExists)

	if !dbExists {
		// Database belum ada — akan dibuat otomatis saat install
		jsonResp(w, map[string]any{
			"ok":   true,
			"note": "Database '" + req.Name + "' belum ada dan akan dibuat otomatis.",
		})
		return
	}

	// Database ada — cek apakah sudah punya tabel (sudah pernah diinstall)
	targetDSN := buildDSN(req.Host, req.Port, req.Name, req.User, req.Password)
	targetDB, err := sql.Open("postgres", targetDSN)
	if err != nil {
		jsonResp(w, map[string]any{"ok": false, "error": "Gagal koneksi ke database: " + err.Error()})
		return
	}
	defer targetDB.Close()

	var tableCount int
	_ = targetDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'`).
		Scan(&tableCount)

	if tableCount > 0 {
		jsonResp(w, map[string]any{
			"ok":    false,
			"error": fmt.Sprintf("Database '%s' sudah memiliki %d tabel. Gunakan database kosong atau ganti nama database.", req.Name, tableCount),
		})
		return
	}

	jsonResp(w, map[string]any{"ok": true, "note": "Database kosong, siap untuk instalasi."})
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type installRequest struct {
	DB       dbConfig `json:"db"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
}

// handleInstallStart menyimpan request dan memulai instalasi di goroutine.
// Langsung return job_id agar JS bisa connect ke SSE stream.
func handleInstallStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResp(w, map[string]any{"ok": false, "error": "request tidak valid"})
		return
	}
	if req.DB.Name == "" || req.Name == "" || req.Email == "" || req.Password == "" {
		jsonResp(w, map[string]any{"ok": false, "error": "semua field wajib diisi"})
		return
	}
	if len(req.Password) < 8 {
		jsonResp(w, map[string]any{"ok": false, "error": "password minimal 8 karakter"})
		return
	}

	jobID, job := newJob()
	go runInstall(req, job)
	jsonResp(w, map[string]any{"ok": true, "job_id": jobID})
}

// handleInstallStream — SSE endpoint untuk progress bar.
func handleInstallStream(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job")
	job, ok := getJob(jobID)
	if !ok {
		http.Error(w, "job tidak ditemukan", 404)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, canFlush := w.(http.Flusher)

	sendEvent := func(ev progressEvent) {
		data, _ := json.Marshal(ev)
		fmt.Fprintf(w, "data: %s\n\n", data)
		if canFlush {
			flusher.Flush()
		}
	}

	for ev := range job.ch {
		sendEvent(ev)
		if ev.Done {
			break
		}
	}
	deleteJob(jobID)
}

// ---------- install runner ----------

func runInstall(req installRequest, job *installJob) {
	emit := func(step string, pct int) {
		job.ch <- progressEvent{Step: step, Percent: pct}
	}
	fail := func(msg string) {
		job.ch <- progressEvent{Done: true, OK: false, Error: msg}
		close(job.ch)
	}

	emit("Menghubungkan ke server PostgreSQL...", 5)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Koneksi ke admin DB
	adminDSN := buildDSN(req.DB.Host, req.DB.Port, "postgres", req.DB.User, req.DB.Password)
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		fail("DSN tidak valid: " + err.Error())
		return
	}
	defer adminDB.Close()
	if err := adminDB.PingContext(ctx); err != nil {
		fail("Tidak dapat terhubung ke PostgreSQL: " + err.Error())
		return
	}

	// Buat database jika belum ada
	var dbExists bool
	_ = adminDB.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, req.DB.Name).
		Scan(&dbExists)

	if !dbExists {
		emit("Membuat database '"+req.DB.Name+"'...", 15)
		if _, err := adminDB.ExecContext(ctx,
			fmt.Sprintf(`CREATE DATABASE "%s"`, req.DB.Name)); err != nil {
			fail("Gagal membuat database: " + err.Error())
			return
		}
	} else {
		emit("Database '"+req.DB.Name+"' ditemukan...", 15)
	}

	// Koneksi ke database target
	targetDSN := buildDSN(req.DB.Host, req.DB.Port, req.DB.Name, req.DB.User, req.DB.Password)
	targetDB, err := sql.Open("postgres", targetDSN)
	if err != nil {
		fail("Gagal koneksi ke database: " + err.Error())
		return
	}
	defer targetDB.Close()
	if err := targetDB.PingContext(ctx); err != nil {
		fail("Gagal koneksi ke database target: " + err.Error())
		return
	}

	// Cek tabel sudah ada
	var tableCount int
	_ = targetDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'`).
		Scan(&tableCount)
	if tableCount > 0 {
		fail(fmt.Sprintf("Database '%s' sudah memiliki tabel. Gunakan database kosong.", req.DB.Name))
		return
	}

	// Baca file migrasi
	entries, err := devmigrations.FS.ReadDir(".")
	if err != nil {
		fail("Gagal baca migrasi: " + err.Error())
		return
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	total := len(files)
	basePct := 20
	migPct := 50 // 20–70% untuk migrasi

	// Jalankan migrasi satu per satu dengan progress
	for i, name := range files {
		pct := basePct + (i * migPct / total)
		emit(fmt.Sprintf("Migrasi %d/%d: %s", i+1, total, name), pct)

		content, err := devmigrations.FS.ReadFile(name)
		if err != nil {
			fail("Gagal baca " + name + ": " + err.Error())
			return
		}
		upSQL := extractUpSection(string(content))
		if upSQL == "" {
			continue
		}
		if _, err := targetDB.ExecContext(ctx, upSQL); err != nil {
			fail("Migrasi " + name + " gagal: " + err.Error())
			return
		}
	}

	emit("Membuat akun superuser...", 75)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fail("Gagal hash password")
		return
	}
	id := uuid.New()
	if _, err := targetDB.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash, role, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, 'superuser', 'active', NOW(), NOW())`,
		id, req.Name, req.Email, string(hash)); err != nil {
		fail("Gagal buat superuser: " + err.Error())
		return
	}

	emit("Menyimpan konfigurasi .env...", 90)

	jwtSecret := generateSecret(32)
	port := "8081"
	if err := writeEnvFile(".env", targetDSN, jwtSecret, port); err != nil {
		fail("Gagal tulis .env: " + err.Error())
		return
	}

	emit("Instalasi selesai!", 100)

	job.ch <- progressEvent{
		Done:   true,
		OK:     true,
		APIURL: "http://localhost:" + port,
		Email:  req.Email,
	}
	close(job.ch)
}

// ---------- helpers ----------

func buildDSN(host, port, name, user, password string) string {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name)
}

func extractUpSection(sqlContent string) string {
	const marker = "-- +migrate Up"
	const endMarker = "-- +migrate Down"
	idx := strings.Index(sqlContent, marker)
	if idx < 0 {
		return strings.TrimSpace(sqlContent)
	}
	start := idx + len(marker)
	end := strings.Index(sqlContent[start:], endMarker)
	if end < 0 {
		return strings.TrimSpace(sqlContent[start:])
	}
	return strings.TrimSpace(sqlContent[start : start+end])
}

func generateSecret(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func writeEnvFile(path, databaseURL, jwtSecret, httpPort string) error {
	existing := map[string]string{}
	if f, err := os.Open(path); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if idx := strings.Index(line, "="); idx > 0 {
				k := strings.TrimSpace(line[:idx])
				v := strings.TrimSpace(line[idx+1:])
				existing[k] = v
			}
		}
		f.Close()
	}
	if existing["JWT_SECRET"] != "" {
		jwtSecret = existing["JWT_SECRET"]
	}
	defaults := map[string]string{
		"APP_NAME":      "vernon-license-api",
		"APP_ENV":       "production",
		"HTTP_PORT":     httpPort,
		"LOG_LEVEL":     "info",
		"JWT_EXP_HOURS": "8",
	}
	for k, v := range defaults {
		if existing[k] == "" {
			existing[k] = v
		}
	}
	existing["DATABASE_URL"] = databaseURL
	existing["JWT_SECRET"] = jwtSecret

	order := []string{"APP_NAME", "APP_ENV", "HTTP_PORT", "LOG_LEVEL",
		"DATABASE_URL", "JWT_SECRET", "JWT_EXP_HOURS"}
	written := map[string]bool{}
	var lines []string
	for _, k := range order {
		if v, ok := existing[k]; ok {
			lines = append(lines, k+"="+v)
			written[k] = true
		}
	}
	for k, v := range existing {
		if !written[k] {
			lines = append(lines, k+"="+v)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0600)
}

func readEnvValue(path, key string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "="); idx > 0 {
			k := strings.TrimSpace(line[:idx])
			v := strings.TrimSpace(line[idx+1:])
			if k == key {
				return v
			}
		}
	}
	return ""
}

func jsonResp(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
