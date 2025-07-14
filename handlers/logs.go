package handlers

import (
	"hospital-system/config"
	"hospital-system/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetLogs(c *fiber.Ctx) error {
	db := config.GetDB()

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	logLevel := c.Query("log_level")
	method := c.Query("method")
	statusCode := c.Query("status_code")

	query := `SELECT id_log, method, path, protocol, status_code, response_time, 
                     user_agent, ip, hostname, body, params, query, email, username, 
                     role, log_level, environment, node_version, pid, timestamp, url, created_at 
              FROM Logs WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if logLevel != "" {
		argCount++
		query += " AND log_level = $" + strconv.Itoa(argCount)
		args = append(args, logLevel)
	}

	if method != "" {
		argCount++
		query += " AND method = $" + strconv.Itoa(argCount)
		args = append(args, method)
	}

	if statusCode != "" {
		argCount++
		query += " AND status_code = $" + strconv.Itoa(argCount)
		args = append(args, statusCode)
	}

	query += " ORDER BY timestamp DESC LIMIT $" + strconv.Itoa(argCount+1) + " OFFSET $" + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener logs: " + err.Error()})
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(
			&log.IDLog, &log.Method, &log.Path, &log.Protocol, &log.StatusCode,
			&log.ResponseTime, &log.UserAgent, &log.IP, &log.Hostname, &log.Body,
			&log.Params, &log.Query, &log.Email, &log.Username, &log.Role,
			&log.LogLevel, &log.Environment, &log.NodeVersion, &log.PID,
			&log.Timestamp, &log.URL, &log.CreatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear log: " + err.Error()})
		}
		logs = append(logs, log)
	}

	return c.JSON(logs)
}

func CreateLog(c *fiber.Ctx) error {
	var log models.Log
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos: " + err.Error()})
	}

	db := config.GetDB()

	err := db.QueryRow(`
        INSERT INTO Logs (method, path, protocol, status_code, response_time, user_agent, 
                         ip, hostname, body, params, query, email, username, role, 
                         log_level, environment, node_version, pid, url) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) 
        RETURNING id_log, timestamp, created_at`,
		log.Method, log.Path, log.Protocol, log.StatusCode, log.ResponseTime,
		log.UserAgent, log.IP, log.Hostname, log.Body, log.Params, log.Query,
		log.Email, log.Username, log.Role, log.LogLevel, log.Environment,
		log.NodeVersion, log.PID, log.URL,
	).Scan(&log.IDLog, &log.Timestamp, &log.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear log: " + err.Error()})
	}

	return c.Status(201).JSON(log)
}

func DeleteOldLogs(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30) 

	db := config.GetDB()

	result, err := db.Exec(`DELETE FROM Logs WHERE created_at < NOW() - INTERVAL '%d days'`, days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar logs: " + err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()

	return c.JSON(fiber.Map{
		"message":       "Logs eliminados exitosamente",
		"deleted_count": rowsAffected,
	})
}
