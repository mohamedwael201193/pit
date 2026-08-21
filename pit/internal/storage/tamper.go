package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func RootsMatch(want, got string) error {
	if err := RejectBadRoot(want); err != nil {
		return err
	}
	if !strings.EqualFold(want, got) {
		return fmt.Errorf("tamper_or_wrong_root")
	}
	return nil
}

func PayloadHash(b []byte) string {
	sum := sha256.Sum256(b)
	return "0x" + hex.EncodeToString(sum[:])
}

func ComparePlain(want, got []byte) error {
	if PayloadHash(want) != PayloadHash(got) {
		return fmt.Errorf("payload_mismatch")
	}
	return nil
}
