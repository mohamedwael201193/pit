package chain8004

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/registry"
)

type Registration struct {
	Network    config.Network
	Owner      identity.Address
	Reporter   identity.Address
	AgentID    string
}

func Prepare(net config.Network, owner, reporter identity.Address) (Registration, error) {
	if owner == "" {
		return Registration{}, fmt.Errorf("owner_required")
	}
	if reporter == "" || reporter == owner {
		return Registration{}, fmt.Errorf("reporter_must_differ")
	}
	addrs := registry.For(net)
	if addrs.Identity8004 == "" || addrs.Reputation8004 == "" {
		return Registration{}, fmt.Errorf("registry_missing")
	}
	return Registration{Network: net, Owner: owner, Reporter: reporter}, nil
}

func FeedbackAllowed(owner, reporter, caller identity.Address) error {
	if caller == owner {
		return fmt.Errorf("owner_self_feedback_reverts")
	}
	if caller != reporter {
		return fmt.Errorf("not_reporter")
	}
	return nil
}

func SameChainIDs(a, b config.Network, idA, idB string) error {
	if a != b {
		return fmt.Errorf("ids_not_portable")
	}
	if strings.TrimSpace(idA) == "" {
		return fmt.Errorf("missing_id")
	}
	_ = idB
	return nil
}
