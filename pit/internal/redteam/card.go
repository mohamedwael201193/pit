package redteam

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func MutatedPreview(before, after engine.Preview) error {
	bh, err := engine.CanonicalHash(before)
	if err != nil {
		return err
	}
	ah, err := engine.CanonicalHash(after)
	if err != nil {
		return err
	}
	if bh == ah {
		return fmt.Errorf("mutation_not_detected")
	}
	return nil
}
