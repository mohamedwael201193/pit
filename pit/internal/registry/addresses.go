package registry

import "github.com/mohamedwael201193/pit/internal/config"

type Addresses struct {
	Identity8004   string
	Reputation8004 string
	DeskID         string
	Serving        string
	Ledger         string
}

func For(n config.Network) Addresses {
	c := config.For(n)
	return Addresses{
		Identity8004:   c.Identity8004,
		Reputation8004: c.Reputation8004,
		DeskID:         c.DeskID,
		Serving:        c.Serving,
		Ledger:         c.Ledger,
	}
}
