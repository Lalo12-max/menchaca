package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsStrongPassword(password string) bool {
	return ValidatePasswordStrength(password) == nil
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 12 {
		return errors.New("la contraseña debe tener al menos 12 caracteres")
	}

	
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("la contraseña debe contener al menos un número")
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("la contraseña debe contener al menos una letra minúscula")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("la contraseña debe contener al menos una letra mayúscula")
	}

	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
	if !hasSymbol {
		return errors.New("la contraseña debe contener al menos un símbolo especial (!@#$%^&*()_+-=[]{}|;':,.<>?~)")
	}

	return nil
}

func HashPasswordWithValidation(password string) (string, error) {
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	return HashPassword(password)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
