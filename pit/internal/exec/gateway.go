package exec

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/session"
)

type Intent struct {
	Action    string
	Preview   engine.Preview
	Hash      string
	Workspace string
}

func RoundPx(px float64, szDecimals int) (string, error) {
	if px <= 0 || math.IsNaN(px) || math.IsInf(px, 0) {
		return "", fmt.Errorf("bad_px")
	}
	prec := 6 - szDecimals
	if prec < 0 {
		prec = 0
	}
	// 5 significant figures, then tick = 10^(szDecimals-6) style rounding.
	g := strconv.FormatFloat(px, 'g', 5, 64)
	f, err := strconv.ParseFloat(g, 64)
	if err != nil {
		return "", err
	}
	pow := math.Pow(10, float64(prec))
	rounded := math.Round(f*pow) / pow
	return strconv.FormatFloat(rounded, 'f', -1, 64), nil
}

func Prepare(in Intent, nowMs int64, used map[string]struct{}) error {
	if err := session.CheckAction(in.Action); err != nil {
		return err
	}
	if in.Preview.WorkspaceID != "" && in.Workspace != "" && in.Preview.WorkspaceID != in.Workspace {
		return fmt.Errorf("wrong_workspace")
	}
	if in.Action == "order" {
		if err := engine.Authorize(in.Preview, in.Hash, nowMs, used); err != nil {
			return err
		}
	}
	if in.Action == "cancel" && strings.TrimSpace(in.Preview.Cloid) == "" && in.Preview.ForecastID == "" {
		return fmt.Errorf("cancel_needs_cloid")
	}
	return nil
}
