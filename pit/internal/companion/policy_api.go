package companion

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/feasibility"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) policyPublic(consequences []string) map[string]any {
	p := policy.Peek(h.Dir)
	hash, _ := p.Hash()
	pinned := false
	workspace := ""
	pin := ""
	if st, err := cli.Load(h.Dir); err == nil {
		workspace = st.WorkspaceID
		pinned = cli.CheckPinned(h.Dir, st.WorkspaceID, p) == nil
		pin, _ = policy.ReadPin(h.Dir, st.WorkspaceID)
	}
	cards := policy.Cards(p)
	rows := make([]map[string]any, 0, len(cards))
	for _, c := range cards {
		rows = append(rows, map[string]any{"title": c.Title, "value": c.Value, "law": c.Law})
	}
	if consequences == nil {
		consequences = []string{}
	}
	allow, refuse := policy.AllowedRefused(p)
	return map[string]any{
		"ok": true, "pinned": pinned, "hash": hash, "pinHash": pin, "version": p.Version, "workspace": workspace,
		"policy": p, "cards": rows, "consequences": consequences,
		"allowed": allow, "refused": refuse,
		"assets": policy.HostAssets, "clipFloor": policy.ClipFloorUSD, "clipCeil": policy.ClipCeilUSD,
		"leverageLocked": 1, "venues": []string{"hyperliquid"}, "marketTypes": []string{"perp"},
		"mutate": false, "sign": false, "trade": false,
		"note": "Only this computer can pin policy. Chat cannot. Leverage stays 1x. Withdraw stays impossible.",
	}
}

func (h *Hub) sessionAliveNow() bool {
	st, err := cli.Load(h.Dir)
	if err != nil || st.Kill {
		return false
	}
	sf, serr := cli.LoadSession(h.Dir)
	if serr != nil {
		return false
	}
	return session.Alive(sf.Meta().Session(), time.Now().UnixMilli())
}

func (h *Hub) policyPinnedNow() bool {
	st, err := cli.Load(h.Dir)
	if err != nil {
		return false
	}
	return cli.CheckPinned(h.Dir, st.WorkspaceID, policy.Peek(h.Dir)) == nil
}

func (h *Hub) capitalNow() feasibility.Account {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return feasibility.Account{PowerSource: "unbound", Note: "This computer is not bound. Capital is unread."}
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return feasibility.Account{PowerSource: "wrong_network", Note: "Wrong network. Capital stays unread."}
	}
	got, err := hl.New(config.For(net)).Capital(st.Wallet)
	if err != nil {
		return feasibility.Account{PowerSource: "venue_unread", Note: "Hyperliquid did not return capital. PIT will not invent balances."}
	}
	return feasibility.FromCapital(got)
}

func (h *Hub) availableUSD() float64 {
	return h.capitalNow().BuyingPower
}

func (h *Hub) annotateWatch(view watch.PublicView, pol policy.Policy) watch.PublicView {
	acct := h.capitalNow()
	if acct.OpenPositions == 0 {
		acct.OpenPositions = h.openPositionCount()
	}
	return watch.ApplyCapital(view, acct, pol, h.sessionAliveNow(), h.policyPinnedNow())
}

func (h *Hub) opportunityConsequences(draft policy.Policy) []string {
	lines := []string{}
	acct := h.capitalNow()
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return lines
	}
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), policy.Clamp(draft))
	if lerr != nil {
		return lines
	}
	exec, research, tight := []string{}, []string{}, []string{}
	sess, pin := h.sessionAliveNow(), true
	for _, c := range cands {
		f := feasibility.FitBook(c.Book, policy.Clamp(draft), acct, sess, pin)
		if f.ExecutionFeasible {
			exec = append(exec, c.Coin)
		} else if f.PolicyEligible || f.ResearchEligible {
			research = append(research, c.Coin)
		}
		if f.Gate == "policy_clip_tight" {
			tight = append(tight, c.Coin)
		}
	}
	if len(exec) > 12 {
		exec = exec[:12]
	}
	if len(research) > 12 {
		research = research[:12]
	}
	if len(exec) == 0 {
		if len(tight) > 0 {
			lines = append(lines, "This account can clear those books' Hyperliquid floors. This clip cannot: "+strings.Join(tight, ", ")+". Raise max trade, preview, then pin. PIT will not invent size.")
		} else if !h.policyPinnedNow() {
			lines = append(lines, "No current live book is execution-feasible under this draft and this account. Research will not spend 0G until you pin host law. PIT will not invent size.")
		} else {
			lines = append(lines, "No current live book is execution-feasible under this draft and this account. PIT will not invent size.")
		}
	} else {
		lines = append(lines, "Would become executable with this account: "+strings.Join(exec, ", ")+".")
	}
	if len(research) > 0 {
		lines = append(lines, "Would stay research-only or blocked: "+strings.Join(research, ", ")+".")
	}
	if acct.Note != "" {
		lines = append(lines, acct.Note)
	}
	return lines
}

func decodePolicyBody(r *http.Request, fallback policy.Policy) policy.Policy {
	if r == nil || r.Body == nil {
		return fallback
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		return fallback
	}
	var body policy.Policy
	_ = json.Unmarshal(raw, &body)
	var wrap struct {
		Policy policy.Policy `json:"policy"`
	}
	_ = json.Unmarshal(raw, &wrap)
	if wrap.Policy.MaxClipUSD > 0 {
		return wrap.Policy
	}
	if body.MaxClipUSD > 0 {
		return body
	}
	return fallback
}
