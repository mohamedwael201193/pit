package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AuthFile struct {
	Provider      string `json:"provider"`
	URL           string `json:"url"`
	Model         string `json:"model"`
	TeeSigner     string `json:"teeSigner"`
	Verifiability string `json:"verifiability"`
	Authorization string `json:"authorization"`
}

func RefuseRouterKey(auth string) error {
	low := strings.ToLower(auth)
	if strings.Contains(low, "app-sk-") {
		return nil
	}
	if strings.Contains(low, "sk-") || strings.Contains(low, "mk-") {
		return fmt.Errorf("router_api_key_denied")
	}
	return nil
}

func DirectAuthPath() string {
	return strings.TrimSpace(os.Getenv("PIT_DIRECT_AUTH_FILE"))
}

func WriteAuth(path string, sku SKU, authorization string) error {
	if err := DenyRouter(sku.URL); err != nil {
		return err
	}
	if err := RefuseRouterKey(authorization); err != nil {
		return err
	}
	if strings.TrimSpace(authorization) == "" {
		return fmt.Errorf("direct_token_required")
	}
	a := AuthFile{
		Provider:      sku.Provider,
		URL:           sku.URL,
		Model:         sku.Model,
		TeeSigner:     sku.TeeSigner,
		Verifiability: sku.Verifiability,
		Authorization: authorization,
	}
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

func MaterializeAsk(dir string, sku SKU, role Role, envelope []byte, authorization string) (DirectJob, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return DirectJob{}, err
	}
	authPath := filepath.Join(dir, "auth.json")
	promptPath := filepath.Join(dir, "prompt.txt")
	outPath := filepath.Join(dir, string(role)+".json")
	if err := WriteAuth(authPath, sku, authorization); err != nil {
		return DirectJob{}, err
	}
	if err := os.WriteFile(promptPath, envelope, 0o600); err != nil {
		return DirectJob{}, err
	}
	return DirectJob{
		AuthPath:      authPath,
		PromptPath:    promptPath,
		OutPath:       outPath,
		Role:          role,
		ProviderURL:   sku.URL,
		OnchainSigner: sku.TeeSigner,
	}, nil
}
