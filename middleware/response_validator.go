package middleware

import (
	"encoding/json"
	"hospital-system/schemas"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)


func ResponseValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		err := c.Next()
		if err != nil {
			return err
		}

		
		contentType := c.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return nil
		}

		
		responseBody := c.Response().Body()
		if len(responseBody) == 0 {
			return nil
		}

		
		var responseData interface{}
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			log.Printf("W - Error al parsear respuesta JSON: %v", err)
			return nil
		}

		
		path := c.Path()
		statusCode := c.Response().StatusCode()

		var validationErr error

		
		switch {
		case path == "/api/v1/auth/register" && statusCode == 201:
			validationErr = schemas.ValidateRegisterResponse(responseData)
		case path == "/api/v1/auth/login" && statusCode == 200:
			if responseMap, ok := responseData.(map[string]interface{}); ok {
				if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
					validationErr = schemas.ValidateLoginMFAResponse(responseData)
				} else {
					validationErr = schemas.ValidateLoginSuccessResponse(responseData)
				}
			}
		case path == "/api/v1/auth/refresh" && statusCode == 200:
			validationErr = schemas.ValidateRefreshTokenResponse(responseData)
		case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
			validationErr = schemas.ValidateEnableMFAResponse(responseData)
		case statusCode >= 400:
			validationErr = schemas.ValidateErrorResponse(responseData)
		}

		if validationErr != nil {
			log.Printf("F - Respuesta inválida para %s [%d]: %v", path, statusCode, validationErr)
		} else {
			log.Printf("S - Respuesta válida para %s [%d]", path, statusCode)
		}

		return nil
	}
}

func SimpleResponseValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

		go func() {
			contentType := c.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				return
			}

			path := c.Path()
			statusCode := c.Response().StatusCode()
			responseBody := c.Response().Body()

			if len(responseBody) == 0 {
				return
			}

			var responseData interface{}
			if err := json.Unmarshal(responseBody, &responseData); err != nil {
				log.Printf("W - Error al parsear respuesta JSON para validación: %v", err)
				return
			}

			var validationErr error
			switch {
			case path == "/api/v1/auth/register" && statusCode == 201:
				validationErr = schemas.ValidateRegisterResponse(responseData)
			case path == "/api/v1/auth/login" && statusCode == 200:
				if responseMap, ok := responseData.(map[string]interface{}); ok {
					if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
						validationErr = schemas.ValidateLoginMFAResponse(responseData)
					} else {
						validationErr = schemas.ValidateLoginSuccessResponse(responseData)
					}
				}
			case path == "/api/v1/auth/refresh" && statusCode == 200:
				validationErr = schemas.ValidateRefreshTokenResponse(responseData)
			case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
				validationErr = schemas.ValidateEnableMFAResponse(responseData)
			case statusCode >= 400:
				validationErr = schemas.ValidateErrorResponse(responseData)
			}

			if validationErr != nil {
				log.Printf("F - Respuesta inválida para %s [%d]: %v", path, statusCode, validationErr)
			} else {
				log.Printf("S - Respuesta válida para %s [%d]", path, statusCode)
			}
		}()

		return nil
	}
}

func ValidateResponseData(path string, statusCode int, responseData interface{}) error {
	switch {
	case path == "/api/v1/auth/register" && statusCode == 201:
		return schemas.ValidateRegisterResponse(responseData)
	case path == "/api/v1/auth/login" && statusCode == 200:
		if responseMap, ok := responseData.(map[string]interface{}); ok {
			if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
				return schemas.ValidateLoginMFAResponse(responseData)
			} else {
				return schemas.ValidateLoginSuccessResponse(responseData)
			}
		}
	case path == "/api/v1/auth/refresh" && statusCode == 200:
		return schemas.ValidateRefreshTokenResponse(responseData)
	case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
		return schemas.ValidateEnableMFAResponse(responseData)
	case statusCode >= 400:
		return schemas.ValidateErrorResponse(responseData)
	}
	return nil
}
