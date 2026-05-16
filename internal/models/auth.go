package models

type UserStatusResponse struct {
	Success bool       `json:"success"`
	Data    UserStatus `json:"data"`
}

type UserStatus struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Plan       string `json:"plan"`
	ScansUsed  int    `json:"scans_used"`
	ScansLimit int    `json:"scans_limit"`
}
