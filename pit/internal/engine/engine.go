package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

type SizerInput struct {
	MarkPx       float64
	SzDecimals   int
	MaxClipUSD   float64
	RequestedUSD float64
	Side         string
	Coin         string
	AllowedCoins []string
	MaxLeverage  int
	RequestedLev int
	Venue        string
	AllowedVenue string
}

type SizedOrder struct {
	Coin        string  `json:"coin"`
	Side        string  `json:"side"`
	Sz          float64 `json:"sz"`
	NotionalUSD float64 `json:"notionalUsd"`
}

func SizeOrder(in SizerInput) (SizedOrder, error) {
	if math.IsNaN(in.MarkPx) || math.IsInf(in.MarkPx, 0) || in.MarkPx <= 0 {
		return SizedOrder{}, fmt.Errorf("bad_mark")
	}
	if in.MaxClipUSD <= 0 || in.RequestedUSD <= 0 {
		return SizedOrder{}, fmt.Errorf("bad_clip")
	}
	side := strings.ToLower(strings.TrimSpace(in.Side))
	if side != "buy" && side != "sell" {
		return SizedOrder{}, fmt.Errorf("bad_side")
	}
	coin := strings.ToUpper(strings.TrimSpace(in.Coin))
	ok := false
	for _, c := range in.AllowedCoins {
		if strings.EqualFold(c, coin) {
			ok = true
			break
		}
	}
	if !ok {
		return SizedOrder{}, fmt.Errorf("coin_not_allowed")
	}
	if in.AllowedVenue != "" && in.Venue != in.AllowedVenue {
		return SizedOrder{}, fmt.Errorf("wrong_venue")
	}
	if in.RequestedLev > 0 && in.MaxLeverage > 0 && in.RequestedLev > in.MaxLeverage {
		return SizedOrder{}, fmt.Errorf("leverage_above_policy")
	}
	usd := in.RequestedUSD
	if usd > in.MaxClipUSD {
		usd = in.MaxClipUSD
	}
	pow := math.Pow(10, float64(in.SzDecimals))
	sz := math.Floor((usd/in.MarkPx)*pow) / pow
	if sz <= 0 {
		return SizedOrder{}, fmt.Errorf("size_rounds_to_zero")
	}
	notional := sz * in.MarkPx
	if notional+1e-9 < 10 {
		if in.MaxClipUSD+1e-9 < 10 {
			return SizedOrder{}, fmt.Errorf("below_min_notional")
		}
		sz = math.Ceil((10.0/in.MarkPx)*pow) / pow
		notional = sz * in.MarkPx
		tick := in.MarkPx / pow
		if notional+1e-9 < 10 || notional > in.MaxClipUSD+tick+1e-9 {
			return SizedOrder{}, fmt.Errorf("below_min_notional")
		}
	}
	if notional > in.MaxClipUSD+in.MarkPx/pow+1e-9 {
		return SizedOrder{}, fmt.Errorf("notional_exceeds_clip")
	}
	return SizedOrder{Coin: coin, Side: side, Sz: sz, NotionalUSD: notional}, nil
}

type RoleJSON struct {
	Survives     *bool  `json:"survives"`
	Kill         *bool  `json:"kill"`
	ProposedSide string `json:"proposed_side"`
}

type Result struct {
	Eligible          bool        `json:"eligible"`
	Deny              string      `json:"deny,omitempty"`
	Size              *SizedOrder `json:"size,omitempty"`
	ChallengerSurvive bool        `json:"challengerSurvive"`
	Reasons           []string    `json:"reasons"`
}

func parseRole(raw json.RawMessage) RoleJSON {
	var r RoleJSON
	_ = json.Unmarshal(raw, &r)
	return r
}

func Evaluate(markPx float64, szDecimals int, maxClip, requested float64, side, coin string, allowed []string, lev, maxLev int, kill bool, researcher, challenger json.RawMessage) Result {
	return EvaluateCommittee(markPx, szDecimals, maxClip, requested, side, coin, allowed, lev, maxLev, kill, researcher, challenger, nil)
}

func EvaluateCommittee(markPx float64, szDecimals int, maxClip, requested float64, side, coin string, allowed []string, lev, maxLev int, kill bool, researcher, challenger, risk json.RawMessage) Result {
	reasons := []string{}
	ch := parseRole(challenger)
	rk := parseRole(risk)
	survive := true
	if ch.Survives != nil && !*ch.Survives {
		survive = false
		reasons = append(reasons, "challenger_killed")
	}
	if ch.Kill != nil && *ch.Kill {
		survive = false
		reasons = append(reasons, "challenger_killed")
	}
	if rk.Kill != nil && *rk.Kill {
		survive = false
		reasons = append(reasons, "risk_killed")
	}
	if kill {
		survive = false
		reasons = append(reasons, "kill_switch")
	}
	rs := parseRole(researcher)
	useSide := side
	if useSide == "" {
		useSide = rs.ProposedSide
	}
	if strings.EqualFold(useSide, "none") || useSide == "" {
		return Result{Eligible: false, Deny: "no_side", Reasons: append(reasons, "no_side")}
	}
	sized, err := SizeOrder(SizerInput{
		MarkPx: markPx, SzDecimals: szDecimals, MaxClipUSD: maxClip, RequestedUSD: requested,
		Side: useSide, Coin: coin, AllowedCoins: allowed, MaxLeverage: maxLev, RequestedLev: lev,
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err != nil {
		return Result{Eligible: false, Deny: err.Error(), Reasons: append(reasons, err.Error())}
	}
	if !survive {
		deny := "challenger_killed"
		for _, r := range reasons {
			if r == "risk_killed" {
				deny = "risk_killed"
				break
			}
		}
		return Result{Eligible: false, Deny: deny, Size: &sized, ChallengerSurvive: false, Reasons: reasons}
	}
	return Result{Eligible: true, Size: &sized, ChallengerSurvive: true, Reasons: reasons}
}

type Preview struct {
	Market        string  `json:"market"`
	Side          string  `json:"side"`
	Sz            float64 `json:"sz"`
	OrderType     string  `json:"orderType"`
	LimitPx       string  `json:"limitPx"`
	SlippageBps   int     `json:"slippageBps"`
	PolicyVersion string  `json:"policyVersion"`
	SessionID     string  `json:"sessionId"`
	ExpiryUnixMs  int64   `json:"expiryUnixMs"`
	Cloid         string  `json:"cloid"`
	ForecastID    string  `json:"forecastId"`
	Nonce         int64   `json:"nonce"`
	WorkspaceID   string  `json:"workspaceId"`
}

func CanonicalHash(p Preview) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func BindPreview(host Preview, model map[string]any) (Preview, error) {
	_ = model // model size/side never win
	if host.Sz <= 0 || host.Market == "" || host.ForecastID == "" || host.Cloid == "" || host.SessionID == "" || host.WorkspaceID == "" {
		return Preview{}, fmt.Errorf("incomplete_preview")
	}
	return host, nil
}

func Authorize(p Preview, presentedHash string, nowMs int64, usedCloids map[string]struct{}) error {
	h, err := CanonicalHash(p)
	if err != nil {
		return err
	}
	if presentedHash == "" || presentedHash != h {
		return fmt.Errorf("preview_hash_mismatch")
	}
	if nowMs >= p.ExpiryUnixMs {
		return fmt.Errorf("preview_expired")
	}
	if _, used := usedCloids[p.Cloid]; used {
		return fmt.Errorf("cloid_replay")
	}
	return nil
}

func PreviewTTL(now time.Time) int64 {
	return now.Add(5 * time.Minute).UnixMilli()
}
