package engine

import "fmt"

func RefuseNoneMarket(m Market) error {
	if m == "none" || m == "" {
		return fmt.Errorf("no_trade")
	}
	return nil
}
