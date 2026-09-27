package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidPatchKeepsData(t *testing.T) {
	oldData := conversations

	conversations = []Conversation{
		{
			ID:            "1",
			CustomerName:  "John Smith",
			CustomerEmail: "john.smith@example.com",
			Subject:       "Unable to reset password",
			Status:        "OPEN",
			Priority:      "HIGH",
			CreatedAt:     "2026-09-10T09:30:00Z",
		},
	}

	t.Cleanup(func() {
		conversations = oldData
	})

	body := strings.NewReader(
		`{"status":"RESOLVED","priority":"URGENT"}`,
	)

	request := httptest.NewRequest(
		http.MethodPatch,
		"/api/conversations/1",
		body,
	)

	response := httptest.NewRecorder()

	conversationByIDHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status code 400, got %d",
			response.Code,
		)
	}

	if conversations[0].Status != "OPEN" {
		t.Errorf(
			"expected status OPEN, got %s",
			conversations[0].Status,
		)
	}

	if conversations[0].Priority != "HIGH" {
		t.Errorf(
			"expected priority HIGH, got %s",
			conversations[0].Priority,
		)
	}
}
