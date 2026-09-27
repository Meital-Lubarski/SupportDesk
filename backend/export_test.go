package main

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConversationsExportHandlerWritesFilteredCSV(t *testing.T) {
	withConversations(t, []Conversation{
		{
			ID:           "1",
			CustomerName: "John Smith",
			Status:       "OPEN",
			Priority:     "HIGH",
			Tags:         []string{"VIP", "Login"},
			AssignedTo:   "Alex Chen",
			CreatedAt:    "2026-09-10T09:30:00Z",
		},
		{
			ID:           "2",
			CustomerName: "Sarah Cohen",
			Status:       "RESOLVED",
			Priority:     "LOW",
			CreatedAt:    "2026-09-11T09:30:00Z",
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/conversations/export?status=OPEN", nil)
	response := httptest.NewRecorder()

	conversationsExportHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.Code)
	}

	if ct := response.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/csv") {
		t.Errorf("expected Content-Type text/csv, got %s", ct)
	}

	if cd := response.Header().Get("Content-Disposition"); !strings.Contains(cd, "conversations.csv") {
		t.Errorf("expected Content-Disposition to name conversations.csv, got %s", cd)
	}

	rows, err := csv.NewReader(response.Body).ReadAll()

	if err != nil {
		t.Fatalf("failed to parse CSV body: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected header row + 1 data row, got %d rows", len(rows))
	}

	if rows[1][1] != "John Smith" || rows[1][7] != "VIP; Login" {
		t.Errorf("unexpected data row: %+v", rows[1])
	}
}

func TestConversationsExportHandlerRejectsNonGet(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/conversations/export", nil)
	response := httptest.NewRecorder()

	conversationsExportHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status code 405, got %d", response.Code)
	}
}
