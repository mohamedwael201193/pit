package exec

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/session"
)

func CoinFromMarket(market string) (string, error) {
	parts := strings.Split(market, ":")
	if len(parts) != 3 || parts[0] != "hyperliquid" || parts[1] != "perp" || parts[2] == "" {
		return "", fmt.Errorf("bad_market")
	}
	return strings.ToUpper(parts[2]), nil
}

func WireFromPreview(p engine.Preview, asset int) (json.RawMessage, error) {
	if err := session.CheckAction("order"); err != nil {
		return nil, err
	}
	if _, err := CoinFromMarket(p.Market); err != nil {
		return nil, err
	}
	buy := strings.EqualFold(p.Side, "buy")
	sz := strconv.FormatFloat(p.Sz, 'f', -1, 64)
	return hl.BuildOrderFlags(asset, buy, p.LimitPx, sz, p.Cloid, p.ReduceOnly)
}

func NeedAssetIndex(ok bool, idx int) error {
	if !ok || idx < 0 {
		return fmt.Errorf("asset_unresolved")
	}
	return nil
}

func RefuseUnsigned(signed bool) error {
	if !signed {
		return fmt.Errorf("exchange_unsigned")
	}
	return nil
}
