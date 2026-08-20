package compute

import (
	"fmt"
	"net/url"
	"strings"
)

type Target struct {
	URL           string
	Model         string
	Provider      string
	TeeSigner     string
	Verifiability string
	Role          string
}

func DenyRouter(rawURL string) error {
	u := strings.ToLower(strings.TrimSpace(rawURL))
	if u == "" {
		return fmt.Errorf("empty_provider_url")
	}
	if strings.Contains(u, "router-api.0g.ai") || strings.Contains(u, "router-api-testnet") {
		return fmt.Errorf("ROUTER_DOWNGRADE_DENIED")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" {
		return fmt.Errorf("bad_provider_url")
	}
	return nil
}

func RequireTeeML(v string) error {
	if !strings.EqualFold(strings.TrimSpace(v), "teeml") {
		return fmt.Errorf("NOT_TEEML")
	}
	return nil
}

func Validate(t Target) error {
	if err := DenyRouter(t.URL); err != nil {
		return err
	}
	if err := RequireTeeML(t.Verifiability); err != nil {
		return err
	}
	if t.TeeSigner == "" || t.Provider == "" || t.Model == "" {
		return fmt.Errorf("incomplete_provider")
	}
	switch t.Role {
	case "researcher", "challenger", "risk":
	default:
		return fmt.Errorf("bad_role")
	}
	return nil
}

// Independence describes what is actually true for a committee run.
type Independence string

const (
	EnvelopeOnly Independence = "prompt_envelope_only_SAME_PROVIDER"
	Providers    Independence = "independent_providers"
	Models       Independence = "independent_models"
)

func Classify(targets []Target) Independence {
	if len(targets) == 0 {
		return EnvelopeOnly
	}
	prov := map[string]struct{}{}
	mod := map[string]struct{}{}
	for _, t := range targets {
		prov[strings.ToLower(t.Provider)] = struct{}{}
		mod[strings.ToLower(t.Model)] = struct{}{}
	}
	if len(prov) > 1 {
		return Providers
	}
	if len(mod) > 1 {
		return Models
	}
	return EnvelopeOnly
}
