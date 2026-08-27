package cli

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/workspace"
)

// Bind attaches YOUR public wallet to this machine. It never stores a key.
// A second wallet cannot overwrite the first without logout --forget.
func Bind(dir, network, wallet string) (DiskState, error) {
	addr, err := identity.NormalizeAddress(wallet)
	if err != nil {
		return DiskState{}, fmt.Errorf("wallet_required")
	}
	net, err := config.ParseNetwork(network)
	if err != nil {
		return DiskState{}, err
	}
	if err := config.RejectGlobalUser(string(addr)); err != nil {
		return DiskState{}, err
	}
	if prev, err := Load(dir); err == nil {
		if err := RefuseNetworkSwitch(prev.Network, string(net)); err != nil {
			return DiskState{}, err
		}
		prevAddr, err := identity.NormalizeAddress(prev.Wallet)
		if err != nil {
			return DiskState{}, err
		}
		if prevAddr != addr {
			return DiskState{}, fmt.Errorf("workspace_owned")
		}
		return prev, nil
	}
	st := workspace.NewStore()
	ws, err := st.Create(addr, net)
	if err != nil {
		return DiskState{}, err
	}
	out := DiskState{
		WorkspaceID: ws.ID,
		Network:     string(ws.Network),
		Wallet:      string(ws.EVM),
	}
	if err := Save(dir, out); err != nil {
		return DiskState{}, err
	}
	return out, nil
}

func PublicBind(st DiskState) map[string]any {
	return map[string]any{
		"ok":        true,
		"workspace": st.WorkspaceID,
		"network":   st.Network,
		"wallet":    st.Wallet,
		"kill":      st.Kill,
		"sign":      false,
		"trade":     false,
	}
}
