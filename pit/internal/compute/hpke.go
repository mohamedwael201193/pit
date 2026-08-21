package compute

import "fmt"

// Official Direct HPKE constants from 0g-pc-e2ee protocol/wire (SPEC §3/§5).
const (
	HPKEKEM  = "0x0020"
	SealInfo = "0g-pc/v1/seal"
	E2EEKey  = "_e2ee"
)

func RequireSealedMessages(fields []string) error {
	for _, f := range fields {
		if f == "messages" {
			return nil
		}
	}
	return fmt.Errorf("prompt_not_sealed")
}

func RefuseCleartextMessages(body map[string]any) error {
	if _, ok := body["messages"]; ok {
		if _, e2 := body[E2EEKey]; !e2 {
			return fmt.Errorf("plaintext_book_denied")
		}
	}
	return nil
}
