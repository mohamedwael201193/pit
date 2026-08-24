package obs

import "fmt"

func Latency(ms int64) error {
	if ms < 0 {
		return fmt.Errorf("bad_latency")
	}
	return nil
}
