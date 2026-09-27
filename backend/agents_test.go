package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withConversations(t *testing.T, data []Conversation) {
	t.Helper()

	oldData := conversations
	conversations = data

	t.Cleanup(func() {
		conversations = oldData
	})
}

func TestAgentsHandlerReturnsRoster(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	response := httptest.NewRecorder()

	agentsHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.Code)
	}

	var got []string

	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != len(agents) {
		t.Errorf("expected %d agents, got %d", len(agents), len(got))
	}
}

func TestConversationsHandlerFiltersByAssignee(t *testing.T) {
	withConversations(t, []Conversation{
		{ID: "1", AssignedTo: "Alex Chen"},
		{ID: "2", AssignedTo: "Priya Patel"},
		{ID: "3", AssignedTo: ""},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/conversations?assignee=Alex+Chen", nil)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	var got []Conversation

	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != 1 || got[0].ID != "1" {
		t.Errorf("expected only conversation 1, got %+v", got)
	}
}

func TestConversationsHandlerFiltersUnassigned(t *testing.T) {
	withConversations(t, []Conversation{
		{ID: "1", AssignedTo: "Alex Chen"},
		{ID: "2", AssignedTo: ""},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/conversations?assignee=unassigned", nil)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	var got []Conversation

	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != 1 || got[0].ID != "2" {
		t.Errorf("expected only conversation 2, got %+v", got)
	}
}

func TestUpdateConversationRejectsUnknownAssignee(t *testing.T) {
	withConversations(t, []Conversation{
		{ID: "1", AssignedTo: "Alex Chen"},
	})

	body := strings.NewReader(`{"assignedTo":"Not A Real Agent"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/conversations/1", body)
	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", response.Code)
	}

	if conversations[0].AssignedTo != "Alex Chen" {
		t.Errorf("expected assignee to stay Alex Chen, got %s", conversations[0].AssignedTo)
	}
}

func TestUpdateConversationCanUnassign(t *testing.T) {
	withConversations(t, []Conversation{
		{ID: "1", AssignedTo: "Alex Chen"},
	})

	body := strings.NewReader(`{"assignedTo":""}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/conversations/1", body)
	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.Code)
	}

	if conversations[0].AssignedTo != "" {
		t.Errorf("expected assignee to be cleared, got %s", conversations[0].AssignedTo)
	}
}
