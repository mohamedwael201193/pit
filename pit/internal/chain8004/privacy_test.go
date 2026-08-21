package chain8004

import "testing"

func TestFeedbackOmitsBook(t *testing.T) {
	_, err := EncodeFeedback(Feedback{AgentID: "0xabc", Score: 71}, map[string]any{"book": "secret"})
	if err == nil {
		t.Fatal("book")
	}
	b, err := EncodeFeedback(Feedback{AgentID: "0xabc", Score: 71, Tags: []string{"resolved"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) == "" {
		t.Fatal("empty")
	}
}
