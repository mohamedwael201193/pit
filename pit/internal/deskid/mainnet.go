package deskid

import "github.com/mohamedwael201193/pit/internal/config"

func RefuseAristotleTransfer() error {
	return RefuseTransfer(config.Mainnet)
}
