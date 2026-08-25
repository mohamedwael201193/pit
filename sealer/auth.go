package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type authFile struct {
	Provider      string `json:"provider"`
	URL           string `json:"url"`
	Model         string `json:"model"`
	TeeSigner     string `json:"teeSigner"`
	Verifiability string `json:"verifiability"`
	Authorization string `json:"authorization"`
}

func loadAuth(path string) (authFile, error) {
	var a authFile
	b, err := os.ReadFile(path)
	if err != nil {
		return a, err
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return a, err
	}
	return a, validateAuth(a)
}

func validateAuth(a authFile) error {
	if err := refuseRouterURL(a.URL); err != nil {
		return err
	}
	if err := refuseRouterKey(a.Authorization); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(a.Verifiability), "teeml") {
		return fmt.Errorf("NOT_TEEML")
	}
	if strings.TrimSpace(a.Provider) == "" || strings.TrimSpace(a.TeeSigner) == "" || strings.TrimSpace(a.Model) == "" {
		return fmt.Errorf("incomplete_provider")
	}
	if strings.TrimSpace(a.Authorization) == "" {
		return fmt.Errorf("direct_token_required")
	}
	return nil
}

func refuseRouterURL(raw string) error {
	u := strings.ToLower(strings.TrimSpace(raw))
	if u == "" {
		return fmt.Errorf("empty_provider_url")
	}
	if strings.Contains(u, "router-api.0g.ai") || strings.Contains(u, "router-api-testnet") {
		return fmt.Errorf("ROUTER_DOWNGRADE_DENIED")
	}
	if !strings.HasPrefix(u, "https://") {
		return fmt.Errorf("bad_provider_url")
	}
	return nil
}

func refuseRouterKey(auth string) error {
	low := strings.ToLower(auth)
	if strings.Contains(low, "app-sk-") {
		return nil
	}
	if strings.Contains(low, "sk-") || strings.Contains(low, "mk-") {
		return fmt.Errorf("router_api_key_denied")
	}
	return nil
}

func requireRole(role string) error {
	switch role {
	case "researcher", "challenger", "risk":
		return nil
	default:
		return fmt.Errorf("bad_role")
	}
}
