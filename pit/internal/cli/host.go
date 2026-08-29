package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
)

func ClipForAccount(pol policy.Policy, availableUSD float64) (float64, error) {
	clip := pol.MaxClipUSD
	if availableUSD >= 0 {
		if availableUSD+1e-9 < policy.ClipFloorUSD {
			return 0, fmt.Errorf("insufficient_margin")
		}
		if availableUSD < clip {
			clip = availableUSD
		}
	}
	return clip, nil
}

func HostPreview(coin, side, forecastID string, book hl.BookSnapshot, pol policy.Policy, sess session.Session, now time.Time, cloid string, nonce int64) (engine.Preview, error) {
	return HostPreviewForAccount(coin, side, forecastID, book, pol, sess, now, cloid, nonce, -1)
}

func HostPreviewForAccount(coin, side, forecastID string, book hl.BookSnapshot, pol policy.Policy, sess session.Session, now time.Time, cloid string, nonce int64, availableUSD float64) (engine.Preview, error) {
	if strings.TrimSpace(forecastID) == "" {
		return engine.Preview{}, fmt.Errorf("forecast_required")
	}
	if sess.ID == "" || sess.Workspace == "" {
		return engine.Preview{}, fmt.Errorf("session_expired")
	}
	if err := hl.MarkFinite(book); err != nil {
		return engine.Preview{}, err
	}
	clip, err := ClipForAccount(pol, availableUSD)
	if err != nil {
		return engine.Preview{}, err
	}
	sized, err := engine.SizeOrder(engine.SizerInput{
		MarkPx:       book.MarkPx,
		SzDecimals:   book.SzDecimals,
		MaxClipUSD:   pol.MaxClipUSD,
		RequestedUSD: clip,
		Side:         side,
		Coin:         coin,
		AllowedCoins: pol.AllowedAssets,
		MaxLeverage:  pol.MaxLeverage,
		RequestedLev: 1,
		Venue:        "hyperliquid",
		AllowedVenue: "hyperliquid",
	})
	if err != nil {
		return engine.Preview{}, err
	}
	p := engine.Preview{
		Market:        "hyperliquid:perp:" + sized.Coin,
		Side:          sized.Side,
		Sz:            sized.Sz,
		OrderType:     "limit",
		LimitPx:       strconv.FormatFloat(book.MarkPx, 'f', -1, 64),
		SlippageBps:   pol.MaxSlippageBps,
		PolicyVersion: pol.Version,
		SessionID:     sess.ID,
		ExpiryUnixMs:  engine.PreviewTTL(now),
		Cloid:         cloid,
		ForecastID:    forecastID,
		Nonce:         nonce,
		WorkspaceID:   sess.Workspace,
	}
	return engine.BindPreview(p, map[string]any{"sz": 99.0, "side": "sell"})
}
