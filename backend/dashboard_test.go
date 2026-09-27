package main

import (
	"testing"
	"time"
)

func TestBuildDashboardSummary(t *testing.T) {
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	source := []Conversation{
		{Status: "OPEN", Priority: "HIGH", Tags: []string{"VIP"}, FollowUpDate: "2026-09-25", AssignedTo: "Alex Chen"},
		{Status: "IN_PROGRESS", Priority: "MEDIUM", Tags: []string{"VIP", "Billing"}, FollowUpDate: "2026-09-30", AssignedTo: "Alex Chen"},
		{Status: "RESOLVED", Priority: "LOW", Tags: []string{"Account"}, FollowUpDate: "2026-09-20", AssignedTo: ""},
	}

	summary := buildDashboardSummary(source, now)

	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}

	if summary.ByStatus["OPEN"] != 1 || summary.ByPriority["HIGH"] != 1 {
		t.Errorf("unexpected status/priority counts: %+v %+v", summary.ByStatus, summary.ByPriority)
	}

	if summary.ByTag["VIP"] != 2 {
		t.Errorf("expected 2 VIP-tagged conversations, got %d", summary.ByTag["VIP"])
	}

	if summary.OverdueFollowUps != 1 {
		t.Errorf("expected 1 overdue follow-up, got %d", summary.OverdueFollowUps)
	}

	if summary.UpcomingFollowUps != 1 {
		t.Errorf("expected 1 upcoming follow-up, got %d", summary.UpcomingFollowUps)
	}

	if summary.ByAgent["Alex Chen"] != 2 {
		t.Errorf("expected 2 conversations assigned to Alex Chen, got %d", summary.ByAgent["Alex Chen"])
	}

	if summary.ByAgent[unassignedLabel] != 1 {
		t.Errorf("expected 1 unassigned conversation, got %d", summary.ByAgent[unassignedLabel])
	}
}
