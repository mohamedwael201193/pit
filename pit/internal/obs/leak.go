package obs

import "fmt"

func RefuseSensitiveLog(msg string) error {
	if Sanitize(msg) == "redacted" {
		return fmt.Errorf("log_leak")
	}
	for _, w := range []string{"private book", "strategy", "positions"} {
		if containsFold(msg, w) {
			return fmt.Errorf("log_leak")
		}
	}
	return nil
}
