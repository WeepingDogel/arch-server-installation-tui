package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// ValidateHostname checks if the given string is a valid hostname.
func ValidateHostname(hostname string) bool {
	if len(hostname) < 1 || len(hostname) > 63 {
		return false
	}
	for _, c := range hostname {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return true
}

// ValidateIP checks if the given string is a valid IPv4 address.
func ValidateIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 || len(p) > 3 {
			return false
		}
		num, err := strconv.Atoi(p)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// ValidatePort checks if the given port number is valid (1-65535).
func ValidatePort(port int) bool {
	return port > 0 && port < 65536
}

// ValidatePassword checks if the password meets minimum requirements.
func ValidatePassword(password string, minLen int) error {
	if len(password) < minLen {
		return fmt.Errorf("password must be at least %d characters", minLen)
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain uppercase, lowercase, and digit characters")
	}
	return nil
}

// ValidateHostPort parses and validates a "host:port" string.
func ValidateHostPort(s string) (string, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid format: expected host:port")
	}
	host := parts[0]
	portStr := parts[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %s", portStr)
	}
	if !ValidatePort(port) {
		return "", 0, fmt.Errorf("port out of range: %d", port)
	}
	return host, port, nil
}

// GeneratePassword generates a cryptographically random password of the given length.
func GeneratePassword(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	if length < 8 {
		length = 16
	}
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}
