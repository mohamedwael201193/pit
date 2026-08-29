package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
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
