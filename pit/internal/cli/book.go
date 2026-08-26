package cli

import (
	"fmt"
	"strings"
)

func ParseAskFlags(args []string) (market, book string, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--market":
			i++
			if i >= len(args) {
				return "", "", fmt.Errorf("empty_envelope")
			}
			market = args[i]
		case "--book":
			i++
			if i >= len(args) {
				return "", "", fmt.Errorf("empty_envelope")
			}
			book = args[i]
		}
	}
	if strings.TrimSpace(market) == "" || strings.TrimSpace(book) == "" {
		return "", "", fmt.Errorf("empty_envelope")
	}
	return market, book, nil
}
