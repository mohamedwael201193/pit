package market

import (
	"fmt"
	"strings"
)

func RequireSource(source, network, symbol string) error {
	if strings.TrimSpace(source) == "" || strings.TrimSpace(network) == "" || strings.TrimSpace(symbol) == "" {
		return fmt.Errorf("quote_unbound")
	}
	return nil
}
