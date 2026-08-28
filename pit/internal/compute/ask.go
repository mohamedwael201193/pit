package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskid"
)

type StageFn func(string)

func notify(stage StageFn, name string) {
	if stage != nil {
		stage(name)
	}
}

// ProductAsk is the sealed private-book entry. Direct fail stops the operation.
func ProductAsk(net config.Network, deskAuthorized bool, bin string) error {
	return ProductAskEnvelope(net, deskAuthorized, bin, nil, nil)
}

func ProductAskEnvelope(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte) error {
	return ProductAskAuth(net, deskAuthorized, bin, publicMarket, privateBook, AuthFile{})
}

func ProductAskAuth(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte, loaded AuthFile) error {
	_, err := ProductAskReport(net, deskAuthorized, bin, publicMarket, privateBook, loaded)
	return err
}

type AskReport struct {
	Roles       []map[string]any `json:"roles"`
	Note        string           `json:"note"`
	Researcher  []byte           `json:"-"`
	Challenger  []byte           `json:"-"`
	Risk        []byte           `json:"-"`
	Deny        string           `json:"deny,omitempty"`
	Eligible    bool             `json:"eligible"`
	Preview     map[string]any   `json:"preview,omitempty"`
	PreviewHash string           `json:"preview_hash,omitempty"`
}

func ProductAskReport(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte, loaded AuthFile) (AskReport, error) {
	return ProductAskReportStage(net, deskAuthorized, bin, publicMarket, privateBook, loaded, "", nil, nil)
}

func ProductAskReportStage(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte, loaded AuthFile, lastPath string, stage StageFn, stop func() bool) (rep AskReport, err error) {
	if err := deskid.BeforeSealedAsk(deskAuthorized); err != nil {
		return AskReport{}, err
	}
	sku := ForNetwork(net)
	if err := SealedAskEnabled(sku); err != nil {
		return AskReport{}, err
	}
	if err := RefuseSKUCopy(net, sku.Model); err != nil {
		return AskReport{}, err
	}
	if err := DenyRouter(sku.URL); err != nil && sku.URL != "" {
		return AskReport{}, err
	}
	if sku.URL == "" {
		return AskReport{}, fmt.Errorf("provider_url_required")
	}
	if err := MustNativeSealer(bin); err != nil {
		return AskReport{}, err
	}
	if strings.TrimSpace(loaded.Authorization) == "" {
		var loadErr error
		loaded, loadErr = LoadEnvAuthFile()
		if loadErr != nil {
			return AskReport{}, loadErr
		}
	}
	if err := RefuseRouterKey(loaded.Authorization); err != nil {
		return AskReport{}, err
	}
	if err := DenyRouter(loaded.URL); err != nil {
		return AskReport{}, err
	}
	if !skuURLMatch(sku.URL, loaded.URL) {
		return AskReport{}, fmt.Errorf("provider_url_mismatch")
	}
	if len(publicMarket) == 0 || len(privateBook) == 0 {
		return AskReport{}, fmt.Errorf("empty_envelope")
	}
	dir, err := os.MkdirTemp("", "pit-ask-")
	if err != nil {
		return AskReport{}, err
	}
	var jobs []DirectJob
	defer func() {
		_ = SavePublicEvidence(lastPath, jobs, err)
		os.RemoveAll(dir)
	}()
	notify(stage, "SEALING_PRIVATE_BOOK")
	envelopes, err := Committee(publicMarket, privateBook)
	if err != nil {
		return AskReport{}, err
	}
	for _, role := range CommitteeRoles() {
		j, merr := MaterializeAsk(dir, sku, role, envelopes[role], loaded.Authorization)
		if merr != nil {
			return AskReport{}, merr
		}
		j.Bin = bin
		jobs = append(jobs, j)
	}
	if stopped(stop) {
		return AskReport{}, fmt.Errorf("research_cancelled")
	}
	notify(stage, "CONTACTING_PRIVATE_PROVIDER")
	if err := RunCommitteeStages(bin, jobs, stage, stop); err != nil {
		rep = AskReport{Note: "Committee stopped. PIT did not fall back to Router.", Roles: publicJobs(jobs)}
		return rep, err
	}
	notify(stage, "DETERMINISTIC_ENGINE")
	rep = AskReport{Note: HonestLabel(IndependenceNote()), Roles: make([]map[string]any, 0, len(jobs))}
	for _, j := range jobs {
		ev := PublicRoleEvidence(j)
		raw := roleJSONFromOut(j.OutPath)
		var parsed map[string]any
		if json.Unmarshal(raw, &parsed) == nil {
			if side, ok := parsed["proposed_side"].(string); ok && side != "" {
				ev["proposed_side"] = side
			}
			if v, ok := parsed["survives"]; ok {
				ev["survives"] = v
			}
			if v, ok := parsed["kill"]; ok {
				ev["kill"] = v
			}
		}
		switch j.Role {
		case Researcher:
			rep.Researcher = raw
		case Challenger:
			rep.Challenger = raw
		case Risk:
			rep.Risk = raw
		}
		rep.Roles = append(rep.Roles, ev)
	}
	return rep, nil
}

func skuURLMatch(skuURL, authURL string) bool {
	a := strings.TrimRight(strings.TrimSpace(skuURL), "/")
	b := strings.TrimRight(strings.TrimSpace(authURL), "/")
	return a != "" && b != "" && strings.EqualFold(a, b)
}

func publicJobs(jobs []DirectJob) []map[string]any {
	out := make([]map[string]any, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, PublicRoleEvidence(j))
	}
	return out
}

func SavePublicEvidence(path string, jobs []DirectJob, runErr error) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	body := map[string]any{"sign": false, "trade": false, "roles": []any{}}
	if runErr != nil {
		body["error"] = runErr.Error()
	}
	roles := make([]any, 0, len(jobs))
	for _, j := range jobs {
		m := PublicRoleEvidence(j)
		b, err := os.ReadFile(j.OutPath)
		if err == nil {
			var extra map[string]any
			if json.Unmarshal(b, &extra) == nil {
				delete(extra, "authorization")
				delete(extra, "prompt")
				delete(extra, "sanitized_output")
				if clip, ok := extra["post_err_clip"].(string); ok {
					extra["post_err_clip"] = redactSecret(clip)
				}
				if verr, ok := extra["verify_err"].(string); ok {
					extra["verify_err"] = redactSecret(verr)
				}
				parsed := roleJSONFromOut(j.OutPath)
				var parsedObj map[string]any
				if json.Unmarshal(parsed, &parsedObj) == nil {
					if side, ok := parsedObj["proposed_side"].(string); ok && side != "" {
						m["proposed_side"] = side
					}
					if v, ok := parsedObj["survives"]; ok {
						m["survives"] = v
					}
					if v, ok := parsedObj["kill"]; ok {
						m["kill"] = v
					}
				}
				for k, v := range extra {
					if _, exists := m[k]; !exists {
						m[k] = v
					}
				}
			}
		}
		roles = append(roles, m)
	}
	body["roles"] = roles
	raw, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return fmt.Errorf("companion_leak")
	}
	return os.WriteFile(path, raw, 0o600)
}

func redactSecret(s string) string {
	low := strings.ToLower(s)
	if strings.Contains(low, "app-sk-") || strings.Contains(low, "bearer ") {
		return "[redacted]"
	}
	return s
}
