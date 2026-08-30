package companion

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/keyring"
)

const memPrefix = "enc:v1:"

func memoryKey(dir string) ([]byte, error) {
	store, err := keyring.OpenProduct(dir)
	if err != nil {
		return nil, err
	}
	got, err := store.Get("memory", "v1")
	if err == nil && len(got) == 32 {
		return got, nil
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := store.Put("memory", "v1", key); err != nil {
		return nil, err
	}
	return key, nil
}

func sealBytes(dir string, plain []byte) (string, error) {
	key, err := memoryKey(dir)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, plain, nil)
	return memPrefix + hex.EncodeToString(out), nil
}

func openBytes(dir, line string) ([]byte, error) {
	raw := strings.TrimSpace(line)
	if !strings.HasPrefix(raw, memPrefix) {
		if looksLikeHexKey(raw) {
			return nil, fmt.Errorf("memory_plaintext_refused")
		}
		return []byte(raw), nil
	}
	key, err := memoryKey(dir)
	if err != nil {
		return nil, err
	}
	bin, err := hex.DecodeString(strings.TrimPrefix(raw, memPrefix))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(bin) < ns {
		return nil, os.ErrInvalid
	}
	return gcm.Open(nil, bin[:ns], bin[ns:], nil)
}

func workingPath(dir string) string {
	return filepath.Join(dir, "memory-working.enc")
}

func writeWorkingMemory(dir string, row map[string]any) {
	if dir == "" || secretful(fmtAny(row)) {
		return
	}
	row["sign"] = false
	row["trade"] = false
	row["updated"] = time.Now().UnixMilli()
	raw, err := json.Marshal(row)
	if err != nil {
		return
	}
	sealed, err := sealBytes(dir, raw)
	if err != nil {
		return
	}
	_ = os.WriteFile(workingPath(dir), []byte(sealed), 0o600)
}

func fmtAny(row map[string]any) string {
	b, _ := json.Marshal(row)
	return string(b)
}

func forgetMemoryFiles(dir string) {
	_ = os.Remove(filepath.Join(dir, "memory-working.json"))
	_ = os.Remove(workingPath(dir))
	_ = os.Remove(experiencePath(dir))
	_ = os.Remove(chatPath(dir))
	_ = os.Remove(threadsPath(dir))
	store, err := keyring.OpenProduct(dir)
	if err == nil {
		_ = store.Delete("memory", "v1")
	}
}
