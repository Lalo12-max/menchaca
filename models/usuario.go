package models

import "time"

// Mantener el enum existente pero agregar el campo RoleID
type TipoUsuario string

const (
	Paciente  TipoUsuario = "paciente"
	Medico    TipoUsuario = "medico"
	Enfermera TipoUsuario = "enfermera"
	Admin     TipoUsuario = "admin"
)

// Nuevo modelo para roles
type Role struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Nuevo modelo para permisos
type Permission struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
}

// Usuario actualizado (manteniendo compatibilidad)
type Usuario struct {
	IDUsuario    int         `json:"id_usuario"`
	Nombre       string      `json:"nombre"`
	Email        string      `json:"email"`
	Password     string      `json:"-"`
	Tipo         TipoUsuario `json:"tipo"`           // Mantener para compatibilidad
	RoleID       *int        `json:"role_id"`        // Nuevo campo opcional
	Role         *Role       `json:"role,omitempty"` // Relación con rol
	RefreshToken string      `json:"-"`
	TokenExpiry  *time.Time  `json:"-"`
	MFAEnabled   bool        `json:"mfa_enabled"`
	MFASecret    string      `json:"-"`
	BackupCodes  []string    `json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type LoginResponse struct {
	AccessToken   string  `json:"access_token"`
	RefreshToken  string  `json:"refresh_token"`
	User          Usuario `json:"user"`
	RequiresMFA   bool    `json:"requires_mfa,omitempty"`
	NeedsMFASetup bool    `json:"needs_mfa_setup,omitempty"` 
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type EnableMFARequest struct {
	Password string `json:"password" validate:"required"`
}

type EnableMFAResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type VerifyMFARequest struct {
	Secret   string `json:"secret" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required"`
}

type DisableMFARequest struct {
	Password string `json:"password" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required"`
}

type RegisterRequest struct {
	Nombre   string      `json:"nombre" validate:"required"`
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=12"`
	Tipo     TipoUsuario `json:"tipo" validate:"required"`
	RoleID   *int        `json:"role_id,omitempty"`
}

type CreateUsuarioRequest struct {
	Nombre   string      `json:"nombre" validate:"required"`
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=12"`
	Tipo     TipoUsuario `json:"tipo" validate:"required"`
	RoleID   *int        `json:"role_id,omitempty"`
}

type InitialMFASetupRequest struct {
	Secret   string `json:"secret" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required"`
}
