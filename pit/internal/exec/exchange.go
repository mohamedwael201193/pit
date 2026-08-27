package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/hl"
)

type Exchange struct {
	URL  string
	HTTP *http.Client
}

func NewExchange(url string) *Exchange {
	return &Exchange{URL: url, HTTP: &http.Client{Timeout: 20 * time.Second}}
}

func (e *Exchange) Guard(raw json.RawMessage) error {
	if e == nil || strings.TrimSpace(e.URL) == "" {
		return fmt.Errorf("missing_exchange")
	}
	if strings.Contains(strings.ToLower(e.URL), "mock") {
		return fmt.Errorf("mock_exchange_denied")
	}
	return hl.AssertActionType(raw)
}

func (e *Exchange) postBytes(raw []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, e.URL, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "pit/1.0")
	resp, err := e.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (e *Exchange) Post(raw json.RawMessage) ([]byte, error) {
	if err := e.Guard(raw); err != nil {
		return nil, err
	}
	return e.postBytes(raw)
}

func (e *Exchange) PostEnvelope(env hl.Envelope) ([]byte, error) {
	if err := e.Guard(env.Action); err != nil {
		return nil, err
	}
	if err := RefuseUnsigned(env.Signed()); err != nil {
		return nil, err
	}
	raw, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	return e.postBytes(raw)
}
