package companion

import (
	"strings"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func venueTradeLink(network, coin string) string {
	base := "https://app.hyperliquid.xyz"
	if strings.EqualFold(strings.TrimSpace(network), "testnet") {
		base = "https://app.hyperliquid-testnet.xyz"
	}
	c := marketCoin(coin)
	if c == "" {
		return base
	}
	return base + "/trade/" + c
}

func marketCoin(market string) string {
	if got, err := pitexec.CoinFromMarket(market); err == nil {
		return got
	}
	return strings.ToUpper(strings.TrimSpace(market))
}

func workspaceNetwork(dir string) string {
	st, err := cli.Load(dir)
	if err != nil || strings.TrimSpace(st.Network) == "" {
		return "mainnet"
	}
	return st.Network
}

func committeeReason(roles []map[string]any) string {
	parts := make([]string, 0, len(roles))
	for _, rm := range roles {
		role := strings.TrimSpace(fmtString(rm["role"]))
		if role == "" {
			continue
		}
		bit := role
		if ver := strings.TrimSpace(fmtString(rm["verify_e2ee"])); ver != "" {
			bit += ":" + ver
		}
		if side := strings.TrimSpace(fmtString(rm["proposed_side"])); side != "" {
			bit += ":" + side
		}
		parts = append(parts, bit)
	}
	return strings.Join(parts, "; ")
}

func (h *Hub) recordPostedOrder(got cli.OrderResult, action, jobID string) {
	if !got.Posted || strings.TrimSpace(got.OID) == "" {
		return
	}
	net := workspaceNetwork(h.Dir)
	coin := marketCoin(got.Market)
	link := venueTradeLink(net, coin)
	status := strings.TrimSpace(got.Status)
	if status == "" {
		status = "posted"
	}
	kind := ""
	switch status {
	case "filled":
		kind = "order.filled"
	case "resting":
		kind = "order.resting"
	case "canceled", "cancelled":
		kind = "order.canceled"
	}
	if kind != "" {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: kind, Market: got.Market,
			Action: action, Status: status, JobID: jobID, PreviewHash: got.Hash,
			OID: got.OID, Link: link, Reason: "oid:" + got.OID,
		})
	}
	h.recordPositionSnapshot(coin, link, got.OID)
}

func (h *Hub) recordPositionSnapshot(coin, link, oid string) {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return
	}
	rows, _, err := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if err != nil {
		return
	}
	want := strings.ToUpper(strings.TrimSpace(coin))
	for _, row := range rows {
		if !strings.EqualFold(row.Coin, want) {
			continue
		}
		sz := strings.TrimSpace(row.Sz)
		if sz == "" || sz == "0" {
			continue
		}
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "position.updated", Market: row.Coin,
			Action: "position", Status: sz, OID: oid, Link: link, Reason: "sz " + sz,
		})
		return
	}
}
