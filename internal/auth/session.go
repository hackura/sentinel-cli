package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeviceSession struct {
	DeviceID        string    `json:"device_id"`
	UserID          string    `json:"user_id"`
	Email           string    `json:"email"`
	Token           string    `json:"token"`
	Status          string    `json:"status"`
	AuthenticatedAt *time.Time `json:"authenticated_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

func PollForToken(ctx context.Context, apiURL, deviceID string) (*DeviceSession, error) {
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()

	timeout := time.After(15 * time.Minute)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("authentication timed out after 15 minutes")
		case <-ticker.C:
			url := fmt.Sprintf("%s/cli/auth/session?device_id=%s", apiURL, deviceID)
			resp, err := client.Get(url)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				continue
			}

			var result struct {
				Success bool          `json:"success"`
				Status  string        `json:"status"`
				Token   string        `json:"token"`
				User    DeviceSession `json:"user"`
				Error   string        `json:"error"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				continue
			}

			if result.Success && result.Status == "success" && result.Token != "" {
				session := result.User
				session.Token = result.Token
				session.DeviceID = deviceID
				return &session, nil
			}

			if result.Status == "expired" {
				return nil, fmt.Errorf("authentication session expired")
			}
		}
	}
}

