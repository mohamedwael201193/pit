package exec

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseVenue(p policy.Policy, venue string) error {
	want := strings.ToLower(strings.TrimSpace(venue))
	for _, v := range p.AllowedVenues {
		if strings.EqualFold(v, want) {
			return nil
		}
	}
	return fmt.Errorf("venue")
}
