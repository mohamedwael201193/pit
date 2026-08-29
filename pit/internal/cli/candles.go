package cli

import (
	"strconv"

	"github.com/mohamedwael201193/pit/internal/skills"
)

func hostCandles(rows []map[string]any) []skills.Candle {
	out := make([]skills.Candle, 0, len(rows))
	for _, row := range rows {
		c := skills.Candle{
			Open:  candleNum(row["open"]),
			High:  candleNum(row["high"]),
			Low:   candleNum(row["low"]),
			Close: candleNum(row["close"]),
		}
		if c.Close == 0 && c.Open == 0 {
			continue
		}
		out = append(out, c)
	}
	return out
}

func candleNum(v any) float64 {
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
