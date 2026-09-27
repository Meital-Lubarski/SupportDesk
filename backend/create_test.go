package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateConversationDefaultsStatusAndPriority(t *testing.T) {
	withConversations(t, []Conversation{{ID: "1"}, {ID: "2"}})

	oldLog := activityLog
	activityLog = map[string][]ActivityEntry{}
	t.Cleanup(func() { activityLog = oldLog })

	body := strings.NewReader(`{
		"customerName": "New Customer",
		"customerEmail": "new@example.com",
		"subject": "Cannot log in"
	}`)

	request := httptest.NewRequest(http.MethodPost, "/api/conversations", body)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status code 201, got %d: %s", response.Code, response.Body.String())
	}

	var created Conversation

	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != "3" {
		t.Errorf("expected next ID to be 3, got %s", created.ID)
	}

	if created.Status != "OPEN" {
		t.Errorf("expected status OPEN, got %s", created.Status)
	}

	if created.Priority != "MEDIUM" {
		t.Errorf("expected default priority MEDIUM, got %s", created.Priority)
	}

	if len(conversations) != 3 {
		t.Errorf("expected 3 conversations after create, got %d", len(conversations))
	}

	if entries := activityLog["3"]; len(entries) != 1 || entries[0].Description != "Conversation created" {
		t.Errorf("expected a 'Conversation created' activity entry, got %+v", entries)
	}
}

func TestCreateConversationRejectsMissingFields(t *testing.T) {
	withConversations(t, []Conversation{})

	body := strings.NewReader(`{"customerName":"","customerEmail":"a@b.com","subject":"Test"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/conversations", body)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", response.Code)
	}

	if len(conversations) != 0 {
		t.Errorf("expected no conversation to be created, got %d", len(conversations))
	}
}

func TestCreateConversationRejectsInvalidEmail(t *testing.T) {
	withConversations(t, []Conversation{})

	body := strings.NewReader(`{"customerName":"A","customerEmail":"not-an-email","subject":"Test"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/conversations", body)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", response.Code)
	}
}

func TestCreateConversationRejectsInvalidPriority(t *testing.T) {
	withConversations(t, []Conversation{})

	body := strings.NewReader(`{
		"customerName": "A",
		"customerEmail": "a@b.com",
		"subject": "Test",
		"priority": "URGENT"
	}`)

	request := httptest.NewRequest(http.MethodPost, "/api/conversations", body)
	response := httptest.NewRecorder()

	conversationsHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", response.Code)
	}
}

func TestNextConversationIDIgnoresGaps(t *testing.T) {
	withConversations(t, []Conversation{{ID: "1"}, {ID: "5"}, {ID: "3"}})

	if got := nextConversationID(); got != "6" {
		t.Errorf("expected next ID to be 6, got %s", got)
	}
}
