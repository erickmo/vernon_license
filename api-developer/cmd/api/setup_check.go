//go:build !wasm

// Package main — setup_check.go menyediakan deteksi kondisi startup sebelum FX dimulai.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// setupIssues menyimpan masalah yang terdeteksi saat startup check.
type setupIssues struct {
	MissingEnv bool   `json:"missing_env"`
	DBError    string `json:"db_error,omitempty"`
	EnvReason  string `json:"env_reason,omitempty"`
}

// checkStartupConditions memeriksa ketersediaan .env dan koneksi DB.
// Returns nil jika semua OK, atau *setupIssues jika ada masalah.
func checkStartupConditions() *setupIssues {
	issues := &setupIssues{}
	hasIssue := false

	// 1. Cek keberadaan .env
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		issues.MissingEnv = true
		issues.EnvReason = "File .env tidak ditemukan di direktori saat ini"
		hasIssue = true
	} else {
		// .env ada, load dulu
		_ = godotenv.Load()
	}

	// 2. Cek DATABASE_URL dan koneksi DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		issues.DBError = "DATABASE_URL tidak diset di environment"
		hasIssue = true
	} else {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			issues.DBError = fmt.Sprintf("gagal membuka koneksi: %v", err)
			hasIssue = true
		} else {
			if err := db.Ping(); err != nil {
				issues.DBError = err.Error()
				hasIssue = true
			}
			_ = db.Close()
		}
	}

	if !hasIssue {
		return nil
	}
	return issues
}

// serveSetupPage menjalankan HTTP server minimal yang menyajikan halaman setup.
// Dipanggil ketika startup conditions tidak terpenuhi.
func serveSetupPage(issues *setupIssues, port string) {
	if port == "" {
		port = "8081"
	}

	// Baca setup.html template
	setupHTML, err := os.ReadFile("web/setup.html")
	if err != nil {
		// Fallback: HTML inline sederhana
		setupHTML = []byte(buildFallbackHTML(issues))
	}

	// Inject issues sebagai JSON ke dalam HTML (replace placeholder)
	issuesJSON, _ := json.Marshal(issues)
	html := strings.Replace(
		string(setupHTML),
		"window.__SETUP_ISSUES__ || {}",
		string(issuesJSON),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(html))
	})
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"setup_required"}`))
	})

	fmt.Printf("\n⚠️  Vernon License membutuhkan konfigurasi.\n")
	fmt.Printf("   Buka http://localhost:%s untuk instruksi setup.\n\n", port)
	if issues.MissingEnv {
		fmt.Println("   ❌ File .env tidak ditemukan")
	}
	if issues.DBError != "" {
		fmt.Printf("   ❌ Database: %s\n", issues.DBError)
	}
	fmt.Println()

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "setup server error: %v\n", err)
		os.Exit(1)
	}
}

// buildFallbackHTML menghasilkan HTML minimal jika web/setup.html tidak ditemukan.
func buildFallbackHTML(issues *setupIssues) string {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html><head><meta charset="UTF-8">
<title>Vernon License — Setup</title>
<style>body{background:#0F0A1A;color:#E2D9F3;font-family:sans-serif;padding:40px;max-width:600px;margin:auto}
h1{color:#4D2975}pre{background:#1A1035;padding:16px;border-radius:8px;color:#26B8B0}
.err{color:#EF4444}.ok{color:#22C55E}</style></head><body>
<h1>🔑 Vernon License — Setup Diperlukan</h1>`)

	if issues.MissingEnv {
		sb.WriteString(`<p class="err">❌ File .env tidak ditemukan.</p>
<p>Jalankan:</p><pre>cp .env.example .env<br># lalu edit DATABASE_URL, JWT_SECRET, PORT</pre>`)
	}
	if issues.DBError != "" {
		sb.WriteString(fmt.Sprintf(`<p class="err">❌ Koneksi database gagal: %s</p>
<p>Pastikan PostgreSQL berjalan dan DATABASE_URL benar di .env</p>`, issues.DBError))
	}
	sb.WriteString(`<p>Setelah fix, jalankan: <code>make dev</code></p>
<button onclick="location.reload()" style="margin-top:16px;padding:10px 20px;background:#4D2975;color:white;border:none;border-radius:6px;cursor:pointer">🔄 Coba Lagi</button>
</body></html>`)
	return sb.String()
}
