package companion

import (
	"strconv"
	"strings"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) policyPublic(consequences []string) map[string]any {
	p := cli.ActivePolicy(h.Dir)
	hash, _ := p.Hash()
	pinned := false
	workspace := ""
	if st, err := cli.Load(h.Dir); err == nil {
		workspace = st.WorkspaceID
		pinned = cli.CheckPinned(h.Dir, st.WorkspaceID, p) == nil
	}
	cards := policy.Cards(p)
	rows := make([]map[string]any, 0, len(cards))
	for _, c := range cards {
		rows = append(rows, map[string]any{"title": c.Title, "value": c.Value, "law": c.Law})
	}
	if consequences == nil {
		consequences = []string{}
	}
	return map[string]any{
		"ok": true, "pinned": pinned, "hash": hash, "version": p.Version, "workspace": workspace,
		"policy": p, "cards": rows, "consequences": consequences,
		"assets": policy.HostAssets, "clipFloor": policy.ClipFloorUSD, "clipCeil": policy.ClipCeilUSD,
		"leverageLocked": 1, "venues": []string{"hyperliquid"}, "marketTypes": []string{"perp"},
		"mutate": false, "sign": false, "trade": false,
		"note": "Only this computer can pin policy. Chat cannot. Leverage stays 1x. Withdraw stays impossible.",
	}
}

func (h *Hub) availableUSD() float64 {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return 0
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return 0
	}
	_, acct, perr := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if perr != nil {
		return 0
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(acct.Withdrawable), 64)
	return v
}

func annotateExec(view watch.PublicView, open int, avail float64, pol policy.Policy) watch.PublicView {
	block, why := policy.ExecWhy(open, avail, pol)
	view.ExecGate = block
	view.ExecWhy = why
	for i := range view.Coins {
		view.Coins[i].ExecGate = block
		view.Coins[i].ExecWhy = why
	}
	if view.Best != nil {
		view.Best.ExecGate = block
		view.Best.ExecWhy = why
	}
	return view
}
