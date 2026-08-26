package compute

import (
	"fmt"
	"os"
	"strings"
)

func LoadEnvelope(marketPath, bookPath string) ([]byte, []byte, error) {
	if strings.TrimSpace(marketPath) == "" || strings.TrimSpace(bookPath) == "" {
		return nil, nil, fmt.Errorf("empty_envelope")
	}
	market, err := os.ReadFile(marketPath)
	if err != nil {
		return nil, nil, fmt.Errorf("empty_envelope")
	}
	book, err := os.ReadFile(bookPath)
	if err != nil {
		return nil, nil, fmt.Errorf("empty_envelope")
	}
	if len(strings.TrimSpace(string(market))) == 0 || len(strings.TrimSpace(string(book))) == 0 {
		return nil, nil, fmt.Errorf("empty_envelope")
	}
	return market, book, nil
}
