package companion

import "testing"

func TestMayHostGuardedExecuteOnlyAutomation(t *testing.T) {
	if !mayHostGuardedExecute("automation") {
		t.Fatal("automation")
	}
	if mayHostGuardedExecute("chat") || mayHostGuardedExecute("research_ui") || mayHostGuardedExecute("") {
		t.Fatal("chat and research_ui must stop at preview")
	}
}
