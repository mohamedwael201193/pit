package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func LiveOnVenue(network, user, cloid string) (bool, error) {
	net, err := config.ParseNetwork(network)
	if err != nil {
		return false, err
	}
	raw, err := hl.New(config.For(net)).OpenOrders(user)
	if err != nil {
		return false, err
	}
	return hl.CloidOnVenue(raw, cloid), nil
}
