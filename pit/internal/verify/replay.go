package verify

import "fmt"

type Used struct {
	hashes map[string]struct{}
}

func NewUsed() *Used {
	return &Used{hashes: map[string]struct{}{}}
}

func (u *Used) File(previewHash string) error {
	if previewHash == "" {
		return fmt.Errorf("bad_preview_hash")
	}
	if _, ok := u.hashes[previewHash]; ok {
		return fmt.Errorf("receipt_replay")
	}
	u.hashes[previewHash] = struct{}{}
	return nil
}
