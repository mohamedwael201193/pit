package policy

import "fmt"

func HaltDaily(p Policy, realizedPnL float64) error {
	if p.DailyLossUSD > 0 && realizedPnL <= -p.DailyLossUSD {
		return fmt.Errorf("daily_loss_halt")
	}
	return nil
}
