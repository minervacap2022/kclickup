package auth

import (
	"fmt"
	"os"
)

func GetToken() (string, error) {
	token := os.Getenv("CLICKUP_API_TOKEN")
	if token == "" {
		return "", fmt.Errorf("CLICKUP_API_TOKEN environment variable is required")
	}
	return token, nil
}

func GetTeamID() string {
	return os.Getenv("CLICKUP_TEAM_ID")
}
