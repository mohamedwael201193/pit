package exec

import (
	"encoding/json"
	"fmt"
)

func ReceiptOID(body []byte) string {
	var top struct {
		Status   string          `json:"status"`
		Response json.RawMessage `json:"response"`
	}
	if json.Unmarshal(body, &top) != nil || top.Status != "ok" {
		return ""
	}
	var resp struct {
		Data struct {
			Statuses []map[string]json.RawMessage `json:"statuses"`
		} `json:"data"`
	}
	if json.Unmarshal(top.Response, &resp) != nil {
		return ""
	}
	for _, st := range resp.Data.Statuses {
		for _, key := range []string{"resting", "filled"} {
			raw, ok := st[key]
			if !ok {
				continue
			}
			var o struct {
				OID any `json:"oid"`
			}
			if json.Unmarshal(raw, &o) != nil {
				continue
			}
			s := fmt.Sprint(o.OID)
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}
