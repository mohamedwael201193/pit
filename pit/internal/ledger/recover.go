package ledger

import "fmt"

const (
	StatusPreviewed = "previewed"
	StatusSigned    = "signed"
	StatusReceipt   = "receipt"
	StatusTimeout   = "timeout"
)

func (s *Store) Mark(workspace, cloid, status, oid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.db.Exec(`UPDATE actions SET status=?, oid=? WHERE cloid=? AND workspace=?`, status, oid, cloid, workspace)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// Recover decides whether a timeout may post again. Never blindly repost.
func Recover(local Record, exchangeOID string, exchangeKnown bool) (repost bool, err error) {
	switch local.Status {
	case StatusPreviewed:
		return true, nil
	case StatusSigned, StatusReceipt:
		if exchangeKnown && exchangeOID != "" {
			return false, nil
		}
		return false, fmt.Errorf("query_exchange_do_not_repost")
	case StatusTimeout:
		if exchangeKnown {
			return false, nil
		}
		return false, fmt.Errorf("query_exchange_do_not_repost")
	default:
		return false, fmt.Errorf("unknown_status")
	}
}
