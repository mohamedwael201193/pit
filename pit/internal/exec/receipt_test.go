package exec

import "testing"

func TestReceiptOIDResting(t *testing.T) {
	body := []byte(`{"status":"ok","response":{"type":"order","data":{"statuses":[{"resting":{"oid":42}}]}}}`)
	if ReceiptOID(body) != "42" {
		t.Fatal(ReceiptOID(body))
	}
}

func TestReceiptOIDKeepsLargeInteger(t *testing.T) {
	body := []byte(`{"status":"ok","response":{"type":"order","data":{"statuses":[{"resting":{"oid":529167222216}}]}}}`)
	if ReceiptOID(body) != "529167222216" {
		t.Fatalf("oid %s", ReceiptOID(body))
	}
}

func TestReceiptOIDFilled(t *testing.T) {
	body := []byte(`{"status":"ok","response":{"type":"order","data":{"statuses":[{"filled":{"oid":"7"}}]}}}`)
	if ReceiptOID(body) != "7" {
		t.Fatal(ReceiptOID(body))
	}
	if ReceiptStatus(body) != "filled" {
		t.Fatal(ReceiptStatus(body))
	}
}

func TestReceiptStatusResting(t *testing.T) {
	body := []byte(`{"status":"ok","response":{"type":"order","data":{"statuses":[{"resting":{"oid":42}}]}}}`)
	if ReceiptStatus(body) != "resting" {
		t.Fatal(ReceiptStatus(body))
	}
}
