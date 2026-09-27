package main

import (
	"strings"
	"testing"
	"time"
)

func TestSuggestPriorityFlagsUrgentKeywords(t *testing.T) {
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	conversation := Conversation{
		Subject:   "Application crashes on login",
		Notes:     "Customer cannot access their account at all.",
		Status:    "OPEN",
		Priority:  "LOW",
		CreatedAt: "2026-09-23T09:00:00Z",
	}

	priority, reason := suggestPriority(conversation, now)

	if priority != "HIGH" {
		t.Errorf("expected HIGH, got %s (reason: %s)", priority, reason)
	}
}

func TestSuggestPriorityLowForGeneralQuestion(t *testing.T) {
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	conversation := Conversation{
		Subject:   "Question about billing",
		Notes:     "Just wondering how invoices are generated.",
		Status:    "OPEN",
		Priority:  "MEDIUM",
		CreatedAt: "2026-09-26T09:00:00Z",
	}

	priority, _ := suggestPriority(conversation, now)

	if priority != "LOW" {
		t.Errorf("expected LOW, got %s", priority)
	}
}

func TestSuggestPriorityBumpsForVIPAndOverdueFollowUp(t *testing.T) {
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	conversation := Conversation{
		Subject:      "Account access restored",
		Status:       "IN_PROGRESS",
		Priority:     "LOW",
		Tags:         []string{"VIP"},
		FollowUpDate: "2026-09-20",
		CreatedAt:    "2026-09-25T09:00:00Z",
	}

	priority, reason := suggestPriority(conversation, now)

	if priority != "HIGH" {
		t.Errorf("expected HIGH, got %s (reason: %s)", priority, reason)
	}

	if !strings.Contains(reason, "VIP") || !strings.Contains(reason, "overdue") {
		t.Errorf("expected reason to mention VIP and overdue, got %q", reason)
	}
}

func TestBuildInsightsMatchesCurrent(t *testing.T) {
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	conversation := Conversation{
		CustomerName: "Daniel Levi",
		Subject:      "Account access restored",
		Status:       "RESOLVED",
		Priority:     "LOW",
		CreatedAt:    "2026-09-12T15:45:00Z",
	}

	insights := buildInsights(conversation, now)

	if !insights.MatchesCurrent {
		t.Errorf("expected suggested priority to match current LOW priority for a resolved conversation")
	}

	if !strings.Contains(insights.Summary, "Daniel Levi") {
		t.Errorf("expected summary to mention the customer name, got %q", insights.Summary)
	}
}
