package cli

import (
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/feasibility"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func LiveAccount(st DiskState) (open int, power float64, err error) {
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return 0, -1, nerr
	}
	c := hl.New(config.For(net))
	got, cerr := c.Capital(st.Wallet)
	if cerr != nil {
		return 0, -1, cerr
	}
	return got.OpenPositions, got.BuyingPower, nil
}

// SizeWatch applies this computer's live capital and pinned policy to a public Watch.
// Website Origin never calls this. CLI and local MCP do, on this machine only.
func SizeWatch(dir string, view watch.PublicView, p policy.Policy) watch.PublicView {
	st, err := Load(dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return view
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return view
	}
	got, cerr := hl.New(config.For(net)).Capital(st.Wallet)
	if cerr != nil {
		return view
	}
	acct := feasibility.FromCapital(got)
	sess := false
	if sf, e := LoadSession(dir); e == nil {
		sess = session.Alive(sf.Meta().Session(), time.Now().UnixMilli())
	}
	pinned := CheckPinned(dir, st.WorkspaceID, p) == nil
	return watch.ApplyCapital(view, acct, p, sess, pinned)
}
