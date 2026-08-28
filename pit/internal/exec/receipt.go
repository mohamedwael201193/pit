package exec

import (
	"encoding/json"
	"strings"
)

func receiptStatuses(body []byte) []map[string]json.RawMessage {
	var top struct {
		Status   string          `json:"status"`
		Response json.RawMessage `json:"response"`
	}
	if json.Unmarshal(body, &top) != nil || top.Status != "ok" {
		return nil
	}
	var resp struct {
		Data struct {
			Statuses []map[string]json.RawMessage `json:"statuses"`
		} `json:"data"`
	}
	if json.Unmarshal(top.Response, &resp) != nil {
		return nil
	}
	return resp.Data.Statuses
}

func ReceiptOID(body []byte) string {
	for _, st := range receiptStatuses(body) {
		for _, key := range []string{"filled", "resting"} {
			raw, ok := st[key]
			if !ok {
				continue
			}
			var o struct {
				OID json.RawMessage `json:"oid"`
			}
			if json.Unmarshal(raw, &o) != nil {
				continue
			}
			s := strings.TrimSpace(string(o.OID))
			s = strings.Trim(s, `"`)
			if s != "" && s != "null" {
				return s
			}
		}
	}
	return ""
}

func ReceiptStatus(body []byte) string {
	for _, st := range receiptStatuses(body) {
		if _, ok := st["filled"]; ok {
			return "filled"
		}
		if _, ok := st["resting"]; ok {
			return "resting"
		}
	}
	return ""
}
