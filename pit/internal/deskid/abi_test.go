package deskid

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCanonicalSelectors(t *testing.T) {
	if Selector("authorizeUsage(uint256,address)") != "0xfa83d14e" {
		t.Fatal(Selector("authorizeUsage(uint256,address)"))
	}
	if Selector("revokeAuthorization(uint256,address)") != "0xc3612ef7" {
		t.Fatal(Selector("revokeAuthorization(uint256,address)"))
	}
	if Selector("ownerOf(uint256)") != "0x6352211e" {
		t.Fatal(Selector("ownerOf(uint256)"))
	}
	got := Selector("isAuthorized(uint256,address)")
	h := crypto.Keccak256([]byte("isAuthorized(uint256,address)"))
	want := "0x" + common.Bytes2Hex(h[:4])
	if got != want {
		t.Fatal(got)
	}
}

func TestEncodeAuthorizeUsagePrefix(t *testing.T) {
	data, err := EncodeAuthorizeUsage(1, "0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52")
	if err != nil {
		t.Fatal(err)
	}
	hex := CalldataHex(data)
	if !strings.HasPrefix(hex, "0xfa83d14e") {
		t.Fatal(hex[:10])
	}
}

func TestEncodeRevokePrefix(t *testing.T) {
	data, err := EncodeRevokeAuthorization(1, "0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(CalldataHex(data), "0xc3612ef7") {
		t.Fatal(CalldataHex(data)[:10])
	}
}

func TestEncodeMintRefusesZeroHash(t *testing.T) {
	if _, err := EncodeMint("0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52", "https://pit0g.vercel.app/", common.Hash{}, "desk"); err == nil {
		t.Fatal("zero hash")
	}
}

func TestEncodeAuthorizeRefusesZeroUser(t *testing.T) {
	if _, err := EncodeAuthorizeUsage(1, "0x0000000000000000000000000000000000000000"); err == nil {
		t.Fatal("zero")
	}
}
