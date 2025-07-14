package handlers

import (
	"fmt"
	"hospital-system/config"
	"hospital-system/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetConsultas(c *fiber.Ctx) error {
	db := config.GetDB()

	rows, err := db.Query(`
        SELECT c.id_consulta, c.consultorio_id, c.medico_id, c.paciente_id, 
               c.tipo, c.horario, c.diagnostico, c.costo,
               m.nombre as medico_nombre,
               p.nombre as paciente_nombre,
               con.nombre as consultorio_nombre
        FROM consultas c
        LEFT JOIN usuarios m ON c.medico_id = m.id_usuario AND m.tipo = 'medico'
        LEFT JOIN usuarios p ON c.paciente_id = p.id_usuario AND p.tipo = 'paciente'
        LEFT JOIN consultorios con ON c.consultorio_id = con.id_consultorio`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultas"})
	}
	defer rows.Close()

	var consultas []map[string]interface{}
	for rows.Next() {
		var consulta models.Consulta
		var medicoNombre, pacienteNombre, consultorioNombre *string

		err := rows.Scan(&consulta.IDConsulta, &consulta.ConsultorioID,
			&consulta.MedicoID, &consulta.PacienteID,
			&consulta.Tipo, &consulta.Horario,
			&consulta.Diagnostico, &consulta.Costo,
			&medicoNombre, &pacienteNombre, &consultorioNombre)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear consulta"})
		}

		consultaCompleta := map[string]interface{}{
			"id_consulta":        consulta.IDConsulta,
			"consultorio_id":     consulta.ConsultorioID,
			"medico_id":          consulta.MedicoID,
			"paciente_id":        consulta.PacienteID,
			"tipo":               consulta.Tipo,
			"horario":            consulta.Horario,
			"diagnostico":        consulta.Diagnostico,
			"costo":              consulta.Costo,
			"medico_nombre":      getValue(medicoNombre),
			"paciente_nombre":    getValue(pacienteNombre),
			"consultorio_nombre": getValue(consultorioNombre),
		}

		consultas = append(consultas, consultaCompleta)
	}

	return c.JSON(consultas)
}

func getValue(s *string) string {
	if s == nil {
		return "N/A"
	}
	return *s
}

func GetConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()
	var consulta models.Consulta

	err = db.QueryRow(`
        SELECT id_consulta, consultorio_id, medico_id, paciente_id, 
               tipo, horario, diagnostico, costo
        FROM consultas WHERE id_consulta = $1`, id).Scan(
		&consulta.IDConsulta, &consulta.ConsultorioID,
		&consulta.MedicoID, &consulta.PacienteID,
		&consulta.Tipo, &consulta.Horario,
		&consulta.Diagnostico, &consulta.Costo)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consulta no encontrada"})
	}

	return c.JSON(consulta)
}

func CreateConsulta(c *fiber.Ctx) error {
	log.Printf("[CONSULTA] 📥 Body recibido: %s", string(c.Body()))
	log.Printf("[CONSULTA] 📋 Content-Type: %s", c.Get("Content-Type"))
	
	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		log.Printf("[CONSULTA] ❌ Error al parsear datos: %v", err)
		log.Printf("[CONSULTA] 📄 Raw body: %s", string(c.Body()))
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Error al parsear datos: %v", err)})
	}

	log.Printf("[CONSULTA] ✅ Datos parseados correctamente:")
	log.Printf("[CONSULTA]   - ConsultorioID: %d", consulta.ConsultorioID)
	log.Printf("[CONSULTA]   - MedicoID: %d", consulta.MedicoID)
	log.Printf("[CONSULTA]   - PacienteID: %d", consulta.PacienteID)
	log.Printf("[CONSULTA]   - Tipo: %s", consulta.Tipo)
	log.Printf("[CONSULTA]   - Horario: %v", consulta.Horario)
	log.Printf("[CONSULTA]   - Diagnostico: %v", consulta.Diagnostico)
	log.Printf("[CONSULTA]   - Costo: %v", consulta.Costo)

	db := config.GetDB()

	log.Printf("[CONSULTA] 💾 Ejecutando INSERT en base de datos...")
	err := db.QueryRow(`
        INSERT INTO consultas (consultorio_id, medico_id, paciente_id, tipo, horario, diagnostico, costo) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_consulta`,
		consulta.ConsultorioID, consulta.MedicoID, consulta.PacienteID,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo).Scan(&consulta.IDConsulta)

	if err != nil {
		log.Printf("[CONSULTA] ❌ Error en base de datos: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Error al crear consulta: %v", err)})
	}

	log.Printf("[CONSULTA] ✅ Consulta creada exitosamente con ID: %d", consulta.IDConsulta)
	return c.Status(201).JSON(consulta)
}

func UpdateConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	_, err = db.Exec(`
        UPDATE consultas 
        SET consultorio_id = $1, medico_id = $2, paciente_id = $3, 
            tipo = $4, horario = $5, diagnostico = $6, costo = $7
        WHERE id_consulta = $8`,
		consulta.ConsultorioID, consulta.MedicoID, consulta.PacienteID,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar consulta"})
	}

	err = db.QueryRow(`
        SELECT id_consulta, consultorio_id, medico_id, paciente_id, 
               tipo, horario, diagnostico, costo
        FROM consultas WHERE id_consulta = $1`, id).Scan(
		&consulta.IDConsulta, &consulta.ConsultorioID,
		&consulta.MedicoID, &consulta.PacienteID,
		&consulta.Tipo, &consulta.Horario,
		&consulta.Diagnostico, &consulta.Costo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consulta actualizada"})
	}

	return c.JSON(consulta)
}

func DeleteConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()

	_, err = db.Exec("DELETE FROM consultas WHERE id_consulta = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar consulta"})
	}

	return c.JSON(fiber.Map{"message": "Consulta eliminada exitosamente"})
}
