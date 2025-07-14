package handlers

import (
	"database/sql"
	"fmt"
	"hospital-system/config"
	"hospital-system/models"
	"hospital-system/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUsuarios(c *fiber.Ctx) error {
	db := config.GetDB()

	fmt.Println("[USUARIOS] 🔍 Iniciando consulta de usuarios")

	rows, err := db.Query(`
		SELECT id_usuario, nombre, email, tipo, role_id, mfa_enabled, created_at, updated_at 
		FROM Usuarios ORDER BY created_at DESC`)
	if err != nil {
		fmt.Printf("[USUARIOS] ❌ Error en consulta SQL: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener usuarios"})
	}
	defer rows.Close()

	fmt.Println("[USUARIOS] ✅ Consulta SQL ejecutada correctamente")

	var usuarios []models.Usuario
	contador := 0
	for rows.Next() {
		var usuario models.Usuario
		var nombre, email sql.NullString
		var tipo sql.NullString
		var roleID sql.NullInt64

		err := rows.Scan(&usuario.IDUsuario, &nombre, &email,
			&tipo, &roleID, &usuario.MFAEnabled,
			&usuario.CreatedAt, &usuario.UpdatedAt)
		if err != nil {
			fmt.Printf("[USUARIOS] ❌ Error al escanear usuario %d: %v\n", contador+1, err)
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear usuario"})
		}

		if nombre.Valid {
			usuario.Nombre = nombre.String
		} else {
			usuario.Nombre = "Sin nombre"
		}

		if email.Valid {
			usuario.Email = email.String
		} else {
			usuario.Email = "Sin email"
		}

		if tipo.Valid {
			usuario.Tipo = models.TipoUsuario(tipo.String)
		} else {
			usuario.Tipo = "paciente" 
		}

		if roleID.Valid {
			roleIDInt := int(roleID.Int64)
			usuario.RoleID = &roleIDInt
		} else {
			usuario.RoleID = nil
		}

		contador++
		fmt.Printf("[USUARIOS] 📄 Usuario %d escaneado: ID=%d, Nombre=%s, Email=%s\n",
			contador, usuario.IDUsuario, usuario.Nombre, usuario.Email)
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("[USUARIOS] ❌ Error durante iteración de filas: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error durante iteración de usuarios"})
	}

	fmt.Printf("[USUARIOS] ✅ Total de usuarios encontrados: %d\n", len(usuarios))
	fmt.Printf("[USUARIOS] 📤 Enviando respuesta JSON con %d usuarios\n", len(usuarios))

	return c.JSON(usuarios)
}

func GetUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()
	var usuario models.Usuario

	err = db.QueryRow(`
		SELECT id_usuario, nombre, email, tipo, role_id, mfa_enabled, created_at, updated_at 
		FROM Usuarios WHERE id_usuario = $1`, id).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo,
		&usuario.RoleID, &usuario.MFAEnabled, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	return c.JSON(usuario)
}

func CreateUsuario(c *fiber.Ctx) error {
	var req models.CreateUsuarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	if !utils.IsStrongPassword(req.Password) {
		return c.Status(400).JSON(fiber.Map{
			"error": "La contraseña debe tener al menos 12 caracteres, incluir mayúsculas, minúsculas, números y símbolos",
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
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
		INSERT INTO Usuarios (nombre, email, password, tipo, role_id, mfa_enabled, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id_usuario, created_at, updated_at`,
		req.Nombre, req.Email, hashedPassword, req.Tipo, roleID, false, time.Now(), time.Now()).Scan(
		&usuario.IDUsuario, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear usuario"})
	}

	usuario.Nombre = req.Nombre
	usuario.Email = req.Email
	usuario.Tipo = req.Tipo
	usuario.RoleID = roleID
	usuario.MFAEnabled = false

	return c.Status(201).JSON(usuario)
}

func UpdateUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req models.CreateUsuarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()
	var usuario models.Usuario

	if req.Password != "" {
		if !utils.IsStrongPassword(req.Password) {
			return c.Status(400).JSON(fiber.Map{
				"error": "La contraseña debe tener al menos 12 caracteres, incluir mayúsculas, minúsculas, números y símbolos",
			})
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
		}

		_, err = db.Exec(`
			UPDATE Usuarios 
			SET nombre = $1, email = $2, password = $3, tipo = $4, role_id = $5, updated_at = CURRENT_TIMESTAMP 
			WHERE id_usuario = $6`,
			req.Nombre, req.Email, hashedPassword, req.Tipo, req.RoleID, id)
	} else {
		_, err = db.Exec(`
			UPDATE Usuarios 
			SET nombre = $1, email = $2, tipo = $3, role_id = $4, updated_at = CURRENT_TIMESTAMP 
			WHERE id_usuario = $5`,
			req.Nombre, req.Email, req.Tipo, req.RoleID, id)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar usuario"})
	}

	err = db.QueryRow(`
		SELECT id_usuario, nombre, email, tipo, role_id, mfa_enabled, created_at, updated_at 
		FROM Usuarios WHERE id_usuario = $1`, id).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo,
		&usuario.RoleID, &usuario.MFAEnabled, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener usuario actualizado"})
	}

	return c.JSON(usuario)
}

func DeleteUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()

	_, err = db.Exec("DELETE FROM Usuarios WHERE id_usuario = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar usuario"})
	}

	return c.JSON(fiber.Map{"message": "Usuario eliminado correctamente"})
}
