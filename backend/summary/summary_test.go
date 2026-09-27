package summary

import (
	"reflect"
	"testing"
)

func TestSummaryConv(t *testing.T) {
	regularImput := []Conversation{
		{Status: "OPEN", Priority: "HIGH"},
		{Status: "OPEN", Priority: "LOW"},
		{Status: "RESOLVED", Priority: "HIGH"},
		{Status: "OPEN", Priority: "HIGH"},
	}
	expected := map[string]int{
		"OPEN_HIGH":     2,
		"OPEN_LOW":      1,
		"RESOLVED_HIGH": 1,
	}
	res1 := SummarizeConversations(regularImput)
	if !reflect.DeepEqual(res1, expected) {
		t.Error("Regular test failed.")
	}

	//Check if it ignores emppty fields
	emptyInput := []Conversation{
		{Status: "OPEN", Priority: ""},
		{Status: "", Priority: "LOW"},
		{Status: "", Priority: ""},
	}

	res2 := SummarizeConversations(emptyInput)
	if len(res2) != 0 {
		t.Error("Empty test failed. expected empty map.")
	}
}
