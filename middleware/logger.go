package middleware

import (
	"encoding/json"
	"hospital-system/config"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		responseTime := int(time.Since(start).Milliseconds())

		method := c.Method()
		path := c.Path()
		protocol := c.Protocol()
		statusCode := c.Response().StatusCode()
		userAgent := c.Get("User-Agent")
		ip := c.IP()
		originalURL := c.OriginalURL()

		var email, username, role *string
		if userEmail := c.Locals("user_email"); userEmail != nil {
			emailStr := userEmail.(string)
			email = &emailStr
		}
		if userName := c.Locals("username"); userName != nil {
			usernameStr := userName.(string)
			username = &usernameStr
		}
		if userRole := c.Locals("user_role"); userRole != nil {
			roleStr := userRole.(string)
			role = &roleStr
		}

		var bodyStr *string
		if len(c.Body()) > 0 {
			body := string(c.Body())
			bodyStr = &body
		}

		var paramsStr *string
		if len(c.AllParams()) > 0 {
			if paramsJSON, marshalErr := json.Marshal(c.AllParams()); marshalErr == nil {
				params := string(paramsJSON)
				paramsStr = &params
			}
		}

		var queryStr *string
		if c.Request().URI().QueryString() != nil {
			query := string(c.Request().URI().QueryString())
			queryStr = &query
		}

		logLevel := "info"
		if statusCode >= 400 && statusCode < 500 {
			logLevel = "warning"
		} else if statusCode >= 500 {
			logLevel = "error"
		}

		hostname, _ := os.Hostname()
		environment := os.Getenv("ENVIRONMENT")
		if environment == "" {
			environment = "development"
		}
		nodeVersion := runtime.Version()
		pid := os.Getpid()

		go func() {
			db := config.GetDB()
			if db != nil {
				_, dbErr := db.Exec(`
                    INSERT INTO Logs (method, path, protocol, status_code, response_time, user_agent, 
                                     ip, hostname, body, params, query, email, username, role, 
                                     log_level, environment, node_version, pid, url) 
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
					method,
					path,
					protocol,
					statusCode,
					responseTime,
					userAgent,
					ip,
					hostname,
					bodyStr,
					paramsStr,
					queryStr,
					email,
					username,
					role,
					logLevel,
					environment,
					nodeVersion,
					pid,
					originalURL,
				)
				if dbErr != nil {
					println("Error logging to database:", dbErr.Error())
				}
			}
		}()

		return err
	}
}
