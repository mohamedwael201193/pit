package watch

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func KillBlocks(p policy.Policy) error {
	if p.KillSwitch {
		return fmt.Errorf("kill_switch")
	}
	return nil
}
