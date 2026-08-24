package market

import (
	"fmt"
	"strings"
)

func RefuseMock(source string) error {
	if strings.EqualFold(strings.TrimSpace(source), "mock") {
		return fmt.Errorf("mock_quote_denied")
	}
	return nil
}
