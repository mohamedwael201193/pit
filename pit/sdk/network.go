package sdk

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/config"
)

func (c Client) BindNetwork(workspace config.Network) error {
	if c.Network != workspace {
		return fmt.Errorf("network_mismatch")
	}
	return nil
}
