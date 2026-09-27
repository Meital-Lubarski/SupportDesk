package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDiffActivityDescribesEachChangedField(t *testing.T) {
	before := Conversation{
		Status:       "OPEN",
		Priority:     "LOW",
		AssignedTo:   "",
		FollowUpDate: "",
		Tags:         []string{"Billing"},
	}

	after := Conversation{
		Status:       "RESOLVED",
		Priority:     "HIGH",
		AssignedTo:   "Alex Chen",
		FollowUpDate: "2026-10-01",
		Tags:         []string{"Billing", "VIP"},
	}

	changes := diffActivity(before, after)

	if len(changes) != 5 {
		t.Fatalf("expected 5 changes, got %d: %+v", len(changes), changes)
	}

	joined := strings.Join(changes, " | ")

	for _, want := range []string{
		"Status changed from open to resolved",
		"Priority changed from low to high",
		"Assigned to Alex Chen",
		"Follow-up date set to 2026-10-01",
		"Tags updated to: Billing, VIP",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected changes to contain %q, got %q", want, joined)
		}
	}
}

func TestDiffActivityIgnoresUnchangedFields(t *testing.T) {
	conversation := Conversation{Status: "OPEN", Priority: "LOW", Tags: []string{"Billing"}}

	changes := diffActivity(conversation, conversation)

	if len(changes) != 0 {
		t.Errorf("expected no changes, got %+v", changes)
	}
}

func TestDescribeAssigneeChangeCoversUnassignAndReassign(t *testing.T) {
	if got := describeAssigneeChange("", "Alex Chen"); got != "Assigned to Alex Chen" {
		t.Errorf("unexpected: %s", got)
	}

	if got := describeAssigneeChange("Alex Chen", ""); got != "Unassigned from Alex Chen" {
		t.Errorf("unexpected: %s", got)
	}

	if got := describeAssigneeChange("Alex Chen", "Priya Patel"); got != "Reassigned from Alex Chen to Priya Patel" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestUpdateConversationRecordsActivity(t *testing.T) {
	withConversations(t, []Conversation{
		{ID: "1", Status: "OPEN", Priority: "LOW", CreatedAt: "2026-09-10T09:30:00Z"},
	})

	oldLog := activityLog
	activityLog = map[string][]ActivityEntry{"1": {{Timestamp: "2026-09-10T09:30:00Z", Description: "Conversation created"}}}
	t.Cleanup(func() { activityLog = oldLog })

	body := strings.NewReader(`{"status":"RESOLVED","priority":"HIGH"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/conversations/1", body)
	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.Code)
	}

	entries := activityLog["1"]

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries (seed + 2 changes), got %d: %+v", len(entries), entries)
	}
}

func TestConversationActivityHandlerReturnsEntries(t *testing.T) {
	withConversations(t, []Conversation{{ID: "1"}})

	oldLog := activityLog
	activityLog = map[string][]ActivityEntry{
		"1": {{Timestamp: "2026-09-10T09:30:00Z", Description: "Conversation created"}},
	}
	t.Cleanup(func() { activityLog = oldLog })

	request := httptest.NewRequest(http.MethodGet, "/api/conversations/1/activity", nil)
	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.Code)
	}

	var entries []ActivityEntry

	if err := json.NewDecoder(response.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(entries) != 1 || entries[0].Description != "Conversation created" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestConversationActivityHandlerRejectsUnknownID(t *testing.T) {
	withConversations(t, []Conversation{{ID: "1"}})

	request := httptest.NewRequest(http.MethodGet, "/api/conversations/999/activity", nil)
	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status code 404, got %d", response.Code)
	}
}
