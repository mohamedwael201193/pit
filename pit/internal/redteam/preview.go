package redteam

import (
	"fmt"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func PreviewMutations(host engine.Preview) []error {
	h, err := engine.CanonicalHash(host)
	if err != nil {
		return []error{err}
	}
	now := time.Now().UnixMilli()
	used := map[string]struct{}{}
	if err := engine.Authorize(host, h, now, used); err != nil {
		return []error{err}
	}
	var errs []error
	mut := host
	mut.Sz = host.Sz + 1
	if err := engine.Authorize(mut, h, now, used); err == nil {
		errs = append(errs, fmt.Errorf("size_mutation_accepted"))
	}
	mut = host
	mut.Side = "sell"
	if host.Side == "sell" {
		mut.Side = "buy"
	}
	if err := engine.Authorize(mut, h, now, used); err == nil {
		errs = append(errs, fmt.Errorf("side_mutation_accepted"))
	}
	if err := engine.Authorize(host, h, host.ExpiryUnixMs+1, used); err == nil {
		errs = append(errs, fmt.Errorf("expired_accepted"))
	}
	return errs
}
