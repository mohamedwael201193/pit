package sdk

import "testing"

func TestRefreshCannotSign(t *testing.T) {
	if BrowserCanSign() || AfterRefreshCanSign() {
		t.Fatal("browser")
	}
	if !DesktopRecoversPreview() {
		t.Fatal("desktop")
	}
}
