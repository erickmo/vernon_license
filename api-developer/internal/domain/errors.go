// Package domain mendefinisikan domain errors untuk Vernon License.
package domain

import "errors"

// Public API errors

// ErrInvalidAPIKey dikembalikan ketika provision API key tidak valid.
var ErrInvalidAPIKey = errors.New("INVALID_API_KEY")

// ErrAlreadyRegistered dikembalikan ketika instance URL sudah terdaftar.
var ErrAlreadyRegistered = errors.New("ALREADY_REGISTERED")

// ErrProductNotFound dikembalikan ketika product slug tidak ditemukan.
var ErrProductNotFound = errors.New("PRODUCT_NOT_FOUND")

// ErrProductInactive dikembalikan ketika product tidak aktif.
var ErrProductInactive = errors.New("PRODUCT_INACTIVE")

// ErrLicenseNotFound dikembalikan ketika license key tidak ditemukan.
var ErrLicenseNotFound = errors.New("LICENSE_NOT_FOUND")

// ErrRateLimitExceeded dikembalikan ketika rate limit terlampaui.
var ErrRateLimitExceeded = errors.New("RATE_LIMIT_EXCEEDED")

// Internal App errors

// ErrAuthInvalidCredentials dikembalikan ketika email atau password salah.
var ErrAuthInvalidCredentials = errors.New("AUTH_INVALID_CREDENTIALS")

// ErrAuthTokenExpired dikembalikan ketika sesi JWT habis.
var ErrAuthTokenExpired = errors.New("AUTH_TOKEN_EXPIRED")

// ErrAuthInsufficientRole dikembalikan ketika user tidak memiliki akses.
var ErrAuthInsufficientRole = errors.New("AUTH_INSUFFICIENT_ROLE")

// ErrCompanyNotFound dikembalikan ketika company tidak ditemukan.
var ErrCompanyNotFound = errors.New("COMPANY_NOT_FOUND")

// ErrProjectNotFound dikembalikan ketika project tidak ditemukan.
var ErrProjectNotFound = errors.New("PROJECT_NOT_FOUND")

// ErrLicenseInvalidTransition dikembalikan ketika perubahan status license tidak valid.
var ErrLicenseInvalidTransition = errors.New("LICENSE_INVALID_TRANSITION")

// ErrProposalNotFound dikembalikan ketika proposal tidak ditemukan.
var ErrProposalNotFound = errors.New("PROPOSAL_NOT_FOUND")

// ErrProposalNotApproved dikembalikan ketika proposal belum disetujui.
var ErrProposalNotApproved = errors.New("PROPOSAL_NOT_APPROVED")

// ErrProposalNotDraft dikembalikan ketika hanya draft yang bisa diedit sales.
var ErrProposalNotDraft = errors.New("PROPOSAL_NOT_DRAFT")

// ErrProposalNotSubmitted dikembalikan ketika hanya submitted yang bisa di-review.
var ErrProposalNotSubmitted = errors.New("PROPOSAL_NOT_SUBMITTED")

// ErrProposalActiveExists dikembalikan ketika masih ada proposal aktif.
var ErrProposalActiveExists = errors.New("PROPOSAL_ACTIVE_EXISTS")

// ErrProposalGenerationFailed dikembalikan ketika PDF proposal gagal dibuat.
var ErrProposalGenerationFailed = errors.New("PROPOSAL_GENERATION_FAILED")

// ErrProductSlugExists dikembalikan ketika slug produk sudah digunakan.
var ErrProductSlugExists = errors.New("PRODUCT_SLUG_EXISTS")

// ErrProductInvalidModule dikembalikan ketika modul tidak tersedia di produk.
var ErrProductInvalidModule = errors.New("PRODUCT_INVALID_MODULE")

// ErrUserEmailExists dikembalikan ketika email sudah digunakan.
var ErrUserEmailExists = errors.New("USER_EMAIL_EXISTS")

// ErrValidationFailed dikembalikan ketika data tidak valid.
var ErrValidationFailed = errors.New("VALIDATION_FAILED")

// ErrUserNotFound dikembalikan ketika user tidak ditemukan.
var ErrUserNotFound = errors.New("USER_NOT_FOUND")
