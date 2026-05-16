package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateDeviceID generates a unique 32-character hex device ID.
func GenerateDeviceID() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// MaskDeviceID masks the device ID to show only the first and last 4 characters.
// Example: a1b2••••••c3d4
func MaskDeviceID(deviceID string) string {
	if len(deviceID) < 8 {
		return deviceID
	}
	return fmt.Sprintf("%s••••••%s", deviceID[:4], deviceID[len(deviceID)-4:])
}
