package ledger

import (
	"sync"
	"testing"
)

func TestExactlyOnceAndIsolation(t *testing.T) {
	dir := t.TempDir()
	a, err := Open(dir, "mainnet", "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := Open(dir, "mainnet", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	ok, err := a.Apply(Record{Workspace: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Cloid: "0x1", Preview: "h", Status: "signed"})
	if err != nil || !ok {
		t.Fatalf("first %v %v", ok, err)
	}
	ok, err = a.Apply(Record{Workspace: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Cloid: "0x1", Preview: "h", Status: "signed"})
	if err != nil || ok {
		t.Fatalf("dup %v %v", ok, err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	res := make([]bool, 2)
	for i := 0; i < 2; i++ {
		i := i
		go func() {
			defer wg.Done()
			ok, err := a.Apply(Record{Workspace: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Cloid: "0x2", Preview: "h2", Status: "signed"})
			if err != nil {
				t.Errorf("concurrent %v", err)
			}
			res[i] = ok
		}()
	}
	wg.Wait()
	wins := 0
	for _, v := range res {
		if v {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("expected one winner, got %v", res)
	}

	ok, err = b.Apply(Record{Workspace: "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", Cloid: "0x1", Preview: "hb", Status: "signed"})
	if err != nil || !ok {
		t.Fatalf("B same cloid in other workspace must apply locally %v %v", ok, err)
	}
	if _, err := a.Get("bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", "0x1"); err == nil {
		t.Fatal("A must not read B")
	}
}
