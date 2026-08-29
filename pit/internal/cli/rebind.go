package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func CommitteeDecisionFromLastResearch(dir, coin string) (compute.AskReport, error) {
	raw, err := os.ReadFile(filepath.Join(dir, "last-research.json"))
	if err != nil {
		return compute.AskReport{}, err
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return compute.AskReport{}, os.ErrInvalid
	}
	var body struct {
		Roles []map[string]any `json:"roles"`
	}
	if json.Unmarshal(raw, &body) != nil || len(body.Roles) < 3 {
		return compute.AskReport{}, os.ErrInvalid
	}
	rep, err := reportFromRoleMaps(body.Roles)
	if err != nil {
		return compute.AskReport{}, err
	}
	st, err := Load(dir)
	if err != nil {
		return rep, err
	}
	p := ActivePolicy(dir)
	_ = CheckPinned(dir, st.WorkspaceID, p)
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = "ETH"
	}
	got := engine.EvaluateCommittee(2500, 4, p.MaxClipUSD, p.MaxClipUSD, "", want, p.AllowedAssets, 1, p.MaxLeverage, p.KillSwitch || st.Kill, rep.Researcher, rep.Challenger, rep.Risk)
	if !got.Eligible {
		deny := got.Deny
		if deny == "" {
			deny = "no_side"
		}
		rep.Deny = deny
		rep.Eligible = false
		rep.Preview = map[string]any{"eligible": false, "deny": deny, "reasons": got.Reasons}
		return rep, nil
	}
	rep.Eligible = true
	if prev, hash, lerr := LoadPreview(dir); lerr == nil && strings.Contains(strings.ToUpper(prev.Market), want) {
		rep.PreviewHash = hash
		rep.Preview = map[string]any{
			"eligible": true, "market": prev.Market, "side": prev.Side, "sz": prev.Sz,
			"orderType": prev.OrderType, "limitPx": prev.LimitPx, "hash": hash, "cloid": prev.Cloid,
			"expiryUnixMs": prev.ExpiryUnixMs,
		}
		return rep, nil
	}
	return rep, nil
}

func ReportFromLastResearch(dir, coin string) (compute.AskReport, error) {
	rep, err := CommitteeDecisionFromLastResearch(dir, coin)
	if err != nil || !rep.Eligible {
		return rep, err
	}
	if strings.TrimSpace(rep.PreviewHash) != "" {
		return rep, nil
	}
	st, err := Load(dir)
	if err != nil {
		return rep, err
	}
	p := ActivePolicy(dir)
	_ = CheckPinned(dir, st.WorkspaceID, p)
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = "ETH"
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return rep, err
	}
	snap, err := hl.New(config.For(net)).PublicBook(want)
	if err != nil || snap.MarkPx <= 0 {
		rep.Deny = "empty_envelope"
		rep.Eligible = false
		rep.Preview = map[string]any{"eligible": false, "deny": "empty_envelope"}
		return rep, nil
	}
	return BindResearchPreview(dir, want, snap, p, st, rep), nil
}

func reportFromRoleMaps(roles []map[string]any) (compute.AskReport, error) {
	if len(roles) < 3 {
		return compute.AskReport{}, os.ErrInvalid
	}
	rep := compute.AskReport{Roles: roles, Note: "last verified committee"}
	ok := 0
	for _, rm := range roles {
		role := strings.ToLower(strings.TrimSpace(fmtString(rm["role"])))
		if strings.EqualFold(fmtString(rm["verify_e2ee"]), "OK") {
			ok++
		}
		b, _ := json.Marshal(map[string]any{
			"proposed_side": rm["proposed_side"],
			"survives":      rm["survives"],
			"kill":          rm["kill"],
		})
		switch role {
		case "researcher":
			rep.Researcher = b
		case "challenger":
			rep.Challenger = b
		case "risk":
			rep.Risk = b
		}
	}
	if ok < 3 {
		return compute.AskReport{}, os.ErrInvalid
	}
	return rep, nil
}

func fmtString(v any) string {
	s, _ := v.(string)
	return s
}
