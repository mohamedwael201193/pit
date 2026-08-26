package cli

import (
	"fmt"
	"strings"
)

func ParsePreviewFlags(args []string) (market, side, forecast string, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--market":
			i++
			if i < len(args) {
				market = args[i]
			}
		case "--side":
			i++
			if i < len(args) {
				side = args[i]
			}
		case "--forecast":
			i++
			if i < len(args) {
				forecast = args[i]
			}
		}
	}
	if strings.TrimSpace(market) == "" || strings.TrimSpace(side) == "" || strings.TrimSpace(forecast) == "" {
		return "", "", "", fmt.Errorf("preview_flags")
	}
	return market, side, forecast, nil
}
