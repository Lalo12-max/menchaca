package handlers

import (
	"database/sql"
	"hospital-system/config"
	"hospital-system/models"
	"hospital-system/utils"
	"strconv" 

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func EnableMFA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	userEmail := c.Locals("user_email").(string)

	var req models.EnableMFARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM Usuarios WHERE id_usuario = $1", userID).Scan(&hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error del servidor"})
	}

	if !utils.CheckPassword(req.Password, hashedPassword) {
		return c.Status(401).JSON(fiber.Map{"error": "Contraseña incorrecta"})
	}

	secret, err := utils.GenerateMFASecret()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar secreto MFA"})
	}

	backupCodes, err := utils.GenerateBackupCodes(8)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar códigos de respaldo"})
	}

	qrCodeURL, err := utils.GenerateQRCode(userEmail, secret, "Hospital System")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar código QR"})
	}

	response := models.EnableMFAResponse{
		Secret:      secret,
		QRCodeURL:   qrCodeURL,
		BackupCodes: backupCodes,
	}

	return c.JSON(response)
}

func VerifyMFA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.VerifyMFARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	if !utils.ValidateTOTP(req.Secret, req.TOTPCode) {
		return c.Status(400).JSON(fiber.Map{"error": "Código TOTP inválido"})
	}

	db := config.GetDB()

	backupCodes, err := utils.GenerateBackupCodes(8)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar códigos de respaldo"})
	}

	_, err = db.Exec(`
		UPDATE Usuarios 
		SET mfa_enabled = TRUE, mfa_secret = $1, backup_codes = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $3`,
		req.Secret, pq.Array(backupCodes), userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al activar MFA"})
	}

	return c.JSON(fiber.Map{
		"message":      "MFA activado correctamente",
		"backup_codes": backupCodes,
	})
}

func DisableMFA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.DisableMFARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	var hashedPassword, mfaSecret string
	var backupCodes pq.StringArray
	err := db.QueryRow(`
		SELECT password, mfa_secret, backup_codes 
		FROM Usuarios WHERE id_usuario = $1 AND mfa_enabled = TRUE`,
		userID).Scan(&hashedPassword, &mfaSecret, &backupCodes)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(fiber.Map{"error": "MFA no está habilitado"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error del servidor"})
	}

	if !utils.CheckPassword(req.Password, hashedPassword) {
		return c.Status(401).JSON(fiber.Map{"error": "Contraseña incorrecta"})
	}

	validCode := false
	if utils.ValidateTOTP(mfaSecret, req.TOTPCode) {
		validCode = true
	} else {
		_, isBackupCode := utils.ValidateBackupCode(backupCodes, req.TOTPCode)
		validCode = isBackupCode
	}

	if !validCode {
		return c.Status(400).JSON(fiber.Map{"error": "Código de verificación inválido"})
	}

	_, err = db.Exec(`
		UPDATE Usuarios 
		SET mfa_enabled = FALSE, mfa_secret = NULL, backup_codes = NULL, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $1`,
		userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al desactivar MFA"})
	}

	return c.JSON(fiber.Map{"message": "MFA desactivado correctamente"})
}

func GetMFAStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	db := config.GetDB()
	var mfaEnabled bool
	err := db.QueryRow("SELECT mfa_enabled FROM Usuarios WHERE id_usuario = $1", userID).Scan(&mfaEnabled)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error del servidor"})
	}

	return c.JSON(fiber.Map{"mfa_enabled": mfaEnabled})
}

func InitialMFASetup(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID de usuario inválido"})
	}

	db := config.GetDB()

	var mfaEnabled bool
	var userEmail string
	err = db.QueryRow("SELECT mfa_enabled, email FROM Usuarios WHERE id_usuario = $1", userID).Scan(&mfaEnabled, &userEmail)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	if mfaEnabled {
		return c.Status(400).JSON(fiber.Map{"error": "El usuario ya tiene MFA configurado"})
	}

	secret, err := utils.GenerateMFASecret()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar secreto MFA"})
	}

	qrCodeURL, err := utils.GenerateQRCode(userEmail, secret, "Hospital System")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar código QR"})
	}

	response := fiber.Map{
		"secret":       secret,
		"qr_code_url":  qrCodeURL,
		"message":      "Escanee el código QR con su aplicación de autenticación",
		"instructions": "1. Abra Google Authenticator o Microsoft Authenticator\n2. Escanee el código QR\n3. Ingrese el código de 6 dígitos para verificar",
	}

	return c.JSON(response)
}

func VerifyInitialMFASetup(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID de usuario inválido"})
	}

	var req models.InitialMFASetupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	var mfaEnabled bool
	err = db.QueryRow("SELECT mfa_enabled FROM Usuarios WHERE id_usuario = $1", userID).Scan(&mfaEnabled)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	if mfaEnabled {
		return c.Status(400).JSON(fiber.Map{"error": "El usuario ya tiene MFA configurado"})
	}

	if !utils.ValidateTOTP(req.Secret, req.TOTPCode) {
		return c.Status(400).JSON(fiber.Map{"error": "Código TOTP inválido"})
	}

	backupCodes, err := utils.GenerateBackupCodes(8)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar códigos de respaldo"})
	}

	_, err = db.Exec(`
		UPDATE Usuarios 
		SET mfa_enabled = TRUE, mfa_secret = $1, backup_codes = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $3`,
		req.Secret, pq.Array(backupCodes), userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al activar MFA"})
	}

	return c.JSON(fiber.Map{
		"message":      "MFA configurado correctamente",
		"backup_codes": backupCodes,
		"instructions": "Guarde estos códigos de respaldo en un lugar seguro. Puede usarlos si pierde acceso a su aplicación de autenticación.",
	})
}
