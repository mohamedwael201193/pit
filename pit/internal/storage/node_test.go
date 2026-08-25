package storage

import "testing"

func TestRefuseNodeClient(t *testing.T) {
	if err := RefuseNodeClient("node_modules/@0glabs/0g-ts-sdk"); err == nil {
		t.Fatal("ts")
	}
	if err := RefuseNodeClient("0g-storage-client"); err != nil {
		t.Fatal(err)
	}
}
