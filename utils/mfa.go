package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)


func GenerateMFASecret() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}


func GenerateQRCode(email, secret, issuer string) (string, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, email, secret, issuer,
	))
	if err != nil {
		return "", err
	}

	
	png, err := qrcode.Encode(key.String(), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("data:image/png;base64,%s",
		base32.StdEncoding.EncodeToString(png)), nil
}


func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}


func GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 4)
		_, err := rand.Read(bytes)
		if err != nil {
			return nil, err
		}
		codes[i] = fmt.Sprintf("%08x", bytes)
	}
	return codes, nil
}


func ValidateBackupCode(codes []string, inputCode string) ([]string, bool) {
	inputCode = strings.ToLower(strings.TrimSpace(inputCode))
	for i, code := range codes {
		if strings.ToLower(code) == inputCode {
			// Remover el código usado
			newCodes := append(codes[:i], codes[i+1:]...)
			return newCodes, true
		}
	}
	return codes, false
}
