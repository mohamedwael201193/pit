package hl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
)

type Client struct {
	InfoURL     string
	ExchangeURL string
	HTTP        *http.Client
}

func New(chain config.Chain) *Client {
	return &Client{
		InfoURL:     chain.HLInfo,
		ExchangeURL: chain.HLExchange,
		HTTP:        &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) postInfo(payload any) (json.RawMessage, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.InfoURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "pit/1.0")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("hl info %d", resp.StatusCode)
	}
	return body, nil
}

type FundingState string

const (
	FundedSpot FundingState = "FUNDED_SPOT"
	FundedPerp FundingState = "FUNDED_PERP"
	Unfunded   FundingState = "UNFUNDED"
)

type AccountView struct {
	Address       string
	PerpValue     float64
	SpotUSDC      float64
	State         FundingState
	WithdrawnNote string
}

func ParseAccount(perpValue, spotUSDC float64) AccountView {
	v := AccountView{PerpValue: perpValue, SpotUSDC: spotUSDC}
	switch {
	case perpValue > 0:
		v.State = FundedPerp
	case spotUSDC > 0:
		v.State = FundedSpot
		v.WithdrawnNote = "Unified USDC lives in spot. This account is funded."
	default:
		v.State = Unfunded
	}
	return v
}

func asFloat(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

func spotUSDCFromClearinghouse(raw json.RawMessage) float64 {
	var st struct {
		Balances []struct {
			Coin  string `json:"coin"`
			Total string `json:"total"`
		} `json:"balances"`
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		return 0
	}
	for _, b := range st.Balances {
		if b.Coin == "USDC" {
			return asFloat(b.Total)
		}
	}
	return 0
}

func perpValueFromClearinghouse(raw json.RawMessage) float64 {
	var st struct {
		MarginSummary struct {
			AccountValue string `json:"accountValue"`
		} `json:"marginSummary"`
	}
	if json.Unmarshal(raw, &st) == nil {
		return asFloat(st.MarginSummary.AccountValue)
	}
	return 0
}

func (c *Client) Account(user string) (AccountView, error) {
	perpRaw, err := c.postInfo(map[string]any{"type": "clearinghouseState", "user": user})
	if err != nil {
		return AccountView{}, err
	}
	spotRaw, err := c.postInfo(map[string]any{"type": "spotClearinghouseState", "user": user})
	if err != nil {
		return AccountView{}, err
	}
	v := ParseAccount(perpValueFromClearinghouse(perpRaw), spotUSDCFromClearinghouse(spotRaw))
	v.Address = user
	return v, nil
}

type BookSnapshot struct {
	Coin         string  `json:"coin"`
	Asset        int     `json:"asset"`
	MarkPx       float64 `json:"markPx"`
	OraclePx     float64 `json:"oraclePx"`
	Funding      float64 `json:"funding"`
	OpenInterest float64 `json:"openInterest"`
	SzDecimals   int     `json:"szDecimals"`
}

func metaCoinDecimals(universe []map[string]any, coin string) int {
	for _, u := range universe {
		name, _ := u["name"].(string)
		if name == coin {
			if d, ok := u["szDecimals"].(float64); ok {
				return int(d)
			}
		}
	}
	return 4
}

type metaCacheEntry struct {
	at       time.Time
	universe []map[string]any
	ctxs     []map[string]any
}

var (
	metaMu    sync.Mutex
	metaCache = map[string]metaCacheEntry{}
)

const metaTTL = 8 * time.Second

func (c *Client) metaAndCtxs() ([]map[string]any, []map[string]any, error) {
	key := c.InfoURL
	metaMu.Lock()
	if e, ok := metaCache[key]; ok && time.Since(e.at) < metaTTL {
		u, x := e.universe, e.ctxs
		metaMu.Unlock()
		return u, x, nil
	}
	metaMu.Unlock()
	raw, err := c.postInfo(map[string]any{"type": "metaAndAssetCtxs"})
	if err != nil {
		return nil, nil, err
	}
	var pack []json.RawMessage
	if err := json.Unmarshal(raw, &pack); err != nil || len(pack) < 2 {
		return nil, nil, fmt.Errorf("meta shape")
	}
	var meta struct {
		Universe []map[string]any `json:"universe"`
	}
	if err := json.Unmarshal(pack[0], &meta); err != nil {
		return nil, nil, err
	}
	var ctxs []map[string]any
	if err := json.Unmarshal(pack[1], &ctxs); err != nil {
		return nil, nil, err
	}
	metaMu.Lock()
	metaCache[key] = metaCacheEntry{at: time.Now(), universe: meta.Universe, ctxs: ctxs}
	metaMu.Unlock()
	return meta.Universe, ctxs, nil
}

func snapshotFromMeta(universe []map[string]any, ctxs []map[string]any, coin string) (BookSnapshot, bool) {
	idx := -1
	for i, u := range universe {
		if name, _ := u["name"].(string); name == coin {
			idx = i
			break
		}
	}
	if idx < 0 || idx >= len(ctxs) {
		return BookSnapshot{}, false
	}
	if _, err := IndexInUniverse(universe, coin); err != nil {
		return BookSnapshot{}, false
	}
	ctx := ctxs[idx]
	return BookSnapshot{
		Coin:         coin,
		Asset:        idx,
		MarkPx:       asFloat(ctx["markPx"]),
		OraclePx:     asFloat(ctx["oraclePx"]),
		Funding:      asFloat(ctx["funding"]),
		OpenInterest: asFloat(ctx["openInterest"]),
		SzDecimals:   metaCoinDecimals(universe, coin),
	}, true
}

func (c *Client) PublicBook(coin string) (BookSnapshot, error) {
	universe, ctxs, err := c.metaAndCtxs()
	if err != nil {
		return BookSnapshot{}, err
	}
	b, ok := snapshotFromMeta(universe, ctxs, coin)
	if !ok {
		return BookSnapshot{}, fmt.Errorf("unknown coin")
	}
	return b, nil
}

// PublicBooks fetches the venue universe once and slices the requested coins.
func (c *Client) PublicBooks(coins []string) ([]BookSnapshot, error) {
	universe, ctxs, err := c.metaAndCtxs()
	if err != nil {
		return nil, err
	}
	out := make([]BookSnapshot, 0, len(coins))
	for _, coin := range coins {
		if b, ok := snapshotFromMeta(universe, ctxs, coin); ok && b.MarkPx > 0 {
			out = append(out, b)
		}
	}
	return out, nil
}
