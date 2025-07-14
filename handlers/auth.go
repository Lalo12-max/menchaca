package handlers

import (
	"database/sql"
	"fmt"
	"hospital-system/config"
	"hospital-system/models"
	"hospital-system/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func sendValidatedResponse(c *fiber.Ctx, statusCode int, data interface{}, validator func(interface{}) error) error {
	if validator != nil {
		if err := validator(data); err != nil {
			log.Printf("Error de validación de respuesta: %v", err)
			return c.Status(500).JSON(utils.NewResponse(500, "E99", fiber.Map{
				"error": "Error interno del servidor",
				"code":  "INTERNAL_ERROR",
			}))
		}
	}

	return c.Status(statusCode).JSON(data)
}

func sendStandardResponse(c *fiber.Ctx, statusCode int, intCode string, data interface{}) error {
	response := utils.NewResponse(statusCode, intCode, data)
	return c.Status(statusCode).JSON(response)
}

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest 
	if err := c.BodyParser(&req); err != nil {
		return sendStandardResponse(c, 400, utils.REGISTER_PARSE_ERROR, fiber.Map{
			"error":       "Error al parsear datos",
			"description": utils.GetCodeDescription(utils.REGISTER_PARSE_ERROR),
		})
	}

	if !utils.IsStrongPassword(req.Password) {
		return sendStandardResponse(c, 400, utils.REGISTER_WEAK_PASSWORD, fiber.Map{
			"error":       "La contraseña debe tener al menos 12 caracteres, incluir mayúsculas, minúsculas, números y símbolos",
			"description": utils.GetCodeDescription(utils.REGISTER_WEAK_PASSWORD),
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_HASH_ERROR, fiber.Map{
			"error":       "Error al procesar contraseña",
			"description": utils.GetCodeDescription(utils.REGISTER_HASH_ERROR),
		})
	}

	mfaSecret, err := utils.GenerateMFASecret()
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_MFA_ERROR, fiber.Map{
			"error":       "Error al generar secreto MFA",
			"description": utils.GetCodeDescription(utils.REGISTER_MFA_ERROR),
		})
	}

	qrCodeURL, err := utils.GenerateQRCode(req.Email, mfaSecret, "Hospital System")
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_MFA_ERROR, fiber.Map{
			"error":       "Error al generar código QR",
			"description": utils.GetCodeDescription(utils.REGISTER_MFA_ERROR),
		})
	}

	db := config.GetDB()
	var usuario models.Usuario

	var roleID *int
	if req.RoleID != nil {
		roleID = req.RoleID
	} else {
		roleMap := map[string]int{
			"admin":     1,
			"medico":    2,
			"enfermera": 3,
			"paciente":  4,
		}
		if id, exists := roleMap[string(req.Tipo)]; exists {
			roleID = &id
		}
	}

	err = db.QueryRow(`
		INSERT INTO Usuarios (nombre, email, password, tipo, role_id, mfa_enabled, mfa_secret, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id_usuario, nombre, email, tipo, created_at, updated_at`,
		req.Nombre, req.Email, hashedPassword, req.Tipo, roleID, true, mfaSecret, time.Now(), time.Now()).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_DB_ERROR, fiber.Map{
			"error":       "Error al crear usuario",
			"description": utils.GetCodeDescription(utils.REGISTER_DB_ERROR),
		})
	}

	responseData := fiber.Map{
		"description": utils.GetCodeDescription(utils.REGISTER_SUCCESS),
		"message":     "Usuario registrado. Configura Microsoft Authenticator con el código QR o secreto.",
		"secret_key":  mfaSecret, 
		"qr_code_url": qrCodeURL, 
	}

	return sendStandardResponse(c, 201, utils.REGISTER_SUCCESS, responseData)
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return sendStandardResponse(c, 400, utils.LOGIN_PARSE_ERROR, fiber.Map{
			"error":       "Error al parsear datos",
			"description": utils.GetCodeDescription(utils.LOGIN_PARSE_ERROR),
		})
	}

	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	db := config.GetDB()
	var usuario models.Usuario
	var hashedPassword string
	var mfaSecret sql.NullString 
	var backupCodes pq.StringArray
	var roleID *int

	err := db.QueryRow(`
		SELECT id_usuario, nombre, email, password, tipo, role_id, mfa_enabled, mfa_secret, backup_codes, created_at, updated_at 
		FROM Usuarios WHERE email = $1`, req.Email).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &hashedPassword, &usuario.Tipo, &roleID,
		&usuario.MFAEnabled, &mfaSecret, &backupCodes, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			logLoginAttempt(db, req.Email, clientIP, userAgent, false)
			return sendStandardResponse(c, 401, utils.LOGIN_INVALID_CREDENTIALS, fiber.Map{
				"error":       "Credenciales inválidas",
				"description": utils.GetCodeDescription(utils.LOGIN_INVALID_CREDENTIALS),
			})
		}
		return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
			"error":       "Error del servidor",
			"description": utils.GetCodeDescription(utils.LOGIN_SERVER_ERROR),
		})
	}

	if !utils.CheckPassword(req.Password, hashedPassword) {
		logLoginAttempt(db, req.Email, clientIP, userAgent, false)
		return sendStandardResponse(c, 401, utils.LOGIN_INVALID_CREDENTIALS, fiber.Map{
			"error":       "Credenciales inválidas",
			"description": utils.GetCodeDescription(utils.LOGIN_INVALID_CREDENTIALS),
		})
	}

	if !usuario.MFAEnabled {
		newMFASecret, err := utils.GenerateMFASecret()
		if err != nil {
			return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
				"error":       "Error al generar MFA",
				"description": "Error interno del servidor",
			})
		}

		newBackupCodes, err := utils.GenerateBackupCodes(8)
		if err != nil {
			return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
				"error":       "Error al generar códigos de respaldo",
				"description": "Error interno del servidor",
			})
		}

		qrCodeURL, err := utils.GenerateQRCode(usuario.Email, newMFASecret, "Hospital System")
		if err != nil {
			return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
				"error":       "Error al generar código QR",
				"description": "Error interno del servidor",
			})
		}

		_, err = db.Exec(`
			UPDATE Usuarios 
			SET mfa_enabled = TRUE, mfa_secret = $1, backup_codes = $2, updated_at = CURRENT_TIMESTAMP 
			WHERE id_usuario = $3`,
			newMFASecret, pq.Array(newBackupCodes), usuario.IDUsuario)

		if err != nil {
			return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
				"error":       "Error al activar MFA",
				"description": "Error interno del servidor",
			})
		}

		usuario.MFAEnabled = true

		responseData := fiber.Map{
			"first_login":    true,
			"mfa_configured": true,
			"secret_key":     newMFASecret,
			"qr_code_url":    qrCodeURL,
			"backup_codes":   newBackupCodes,
			"message":        "¡Bienvenido! Se ha configurado automáticamente tu autenticación de dos factores.",
			"instructions":   "1. Guarda estos códigos de respaldo en un lugar seguro\n2. Escanea el código QR con Google Authenticator o Microsoft Authenticator\n3. En futuros logins necesitarás el código de 6 dígitos de tu aplicación",
			"user":           usuario,
		}
		return sendStandardResponse(c, 200, "MFA_AUTO_CONFIGURED", responseData)
	}

	if usuario.MFAEnabled {
		if req.TOTPCode == "" {
			responseData := fiber.Map{
				"requires_mfa": true,
				"message":      "Se requiere código MFA",
				"description":  utils.GetCodeDescription(utils.LOGIN_MFA_REQUIRED),
			}
			return sendStandardResponse(c, 200, utils.LOGIN_MFA_REQUIRED, responseData)
		}

		validMFA := false
		if mfaSecret.Valid && utils.ValidateTOTP(mfaSecret.String, req.TOTPCode) {
			validMFA = true
		} else {
			newBackupCodes, isBackupCode := utils.ValidateBackupCode(backupCodes, req.TOTPCode)
			if isBackupCode {
				validMFA = true
				_, _ = db.Exec("UPDATE Usuarios SET backup_codes = $1 WHERE id_usuario = $2",
					pq.Array(newBackupCodes), usuario.IDUsuario)
			}
		}

		if !validMFA {
			logLoginAttempt(db, req.Email, clientIP, userAgent, false)
			return sendStandardResponse(c, 401, utils.LOGIN_INVALID_MFA, fiber.Map{
				"error":       "Código MFA inválido",
				"description": utils.GetCodeDescription(utils.LOGIN_INVALID_MFA),
			})
		}
	}

	var permissions []utils.Permission
	if roleID != nil {
		permissionRows, err := db.Query(`
			SELECT p.resource, p.action 
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id = $1
		`, *roleID)
		if err == nil {
			defer permissionRows.Close()
			for permissionRows.Next() {
				var permission utils.Permission
				permissionRows.Scan(&permission.Resource, &permission.Action)
				permissions = append(permissions, permission)
			}
		}
	}

	logLoginAttempt(db, req.Email, clientIP, userAgent, true)

	accessToken, err := utils.GenerateAccessToken(usuario.IDUsuario, usuario.Email, string(usuario.Tipo), roleID, permissions)
	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_TOKEN_ERROR, fiber.Map{
			"error":       "Error al generar token",
			"description": utils.GetCodeDescription(utils.LOGIN_TOKEN_ERROR),
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(usuario.IDUsuario)
	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_TOKEN_ERROR, fiber.Map{
			"error":       "Error al generar refresh token",
			"description": utils.GetCodeDescription(utils.LOGIN_TOKEN_ERROR),
		})
	}

	tokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	_, err = db.Exec(`
		UPDATE Usuarios SET refresh_token = $1, token_expiry = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $3`,
		refreshToken, tokenExpiry, usuario.IDUsuario)

	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
			"error":       "Error al guardar refresh token",
			"description": utils.GetCodeDescription(utils.LOGIN_SERVER_ERROR),
		})
	}

	responseData := fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          usuario,
		"message":       "Login exitoso",
		"description":   utils.GetCodeDescription(utils.LOGIN_SUCCESS),
	}

	return sendStandardResponse(c, 200, utils.LOGIN_SUCCESS, responseData)
}

func RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()
	var usuario models.Usuario
	var tokenExpiry time.Time
	var roleID *int

	err := db.QueryRow(`
		SELECT id_usuario, nombre, email, tipo, role_id, token_expiry 
		FROM Usuarios WHERE refresh_token = $1`, req.RefreshToken).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo, &roleID, &tokenExpiry)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token inválido"})
	}

	if time.Now().After(tokenExpiry) {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token expirado"})
	}

	var permissions []utils.Permission
	if roleID != nil {
		permissionRows, err := db.Query(`
			SELECT p.resource, p.action 
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id = $1
		`, *roleID)
		if err == nil {
			defer permissionRows.Close()
			for permissionRows.Next() {
				var permission utils.Permission
				permissionRows.Scan(&permission.Resource, &permission.Action)
				permissions = append(permissions, permission)
			}
		}
	}

	accessToken, err := utils.GenerateAccessToken(usuario.IDUsuario, usuario.Email, string(usuario.Tipo), roleID, permissions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar token"})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func logLoginAttempt(db *sql.DB, email, ip, userAgent string, success bool) {
	status := "FAILED"
	if success {
		status = "SUCCESS"
	}

	fmt.Printf("[LOGIN %s] Email: %s, IP: %s, Time: %s\n", status, email, ip, time.Now().Format(time.RFC3339))
}
