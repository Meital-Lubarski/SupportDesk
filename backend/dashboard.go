package main

import (
	"net/http"
	"time"
)

// Aggregated counts for the dashboard view.
type DashboardSummary struct {
	Total                  int            `json:"total"`
	ByStatus               map[string]int `json:"byStatus"`
	ByPriority             map[string]int `json:"byPriority"`
	ByTag                  map[string]int `json:"byTag"`
	ByAgent                map[string]int `json:"byAgent"`
	OverdueFollowUps       int            `json:"overdueFollowUps"`
	UpcomingFollowUps      int            `json:"upcomingFollowUps"`
	SuggestedPriorityBumps int            `json:"suggestedPriorityBumps"`
}

const unassignedLabel = "Unassigned"

const upcomingFollowUpWindow = 7 * 24 * time.Hour

// Returns aggregated counts across all conversations.
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, buildDashboardSummary(conversations, time.Now()))
}

func buildDashboardSummary(source []Conversation, now time.Time) DashboardSummary {
	summary := DashboardSummary{
		Total:      len(source),
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
		ByTag:      make(map[string]int),
		ByAgent:    make(map[string]int),
	}

	upcomingCutoff := now.Add(upcomingFollowUpWindow)

	for _, conversation := range source {
		summary.ByStatus[conversation.Status]++
		summary.ByPriority[conversation.Priority]++

		for _, tag := range conversation.Tags {
			summary.ByTag[tag]++
		}

		if conversation.AssignedTo == "" {
			summary.ByAgent[unassignedLabel]++
		} else {
			summary.ByAgent[conversation.AssignedTo]++
		}

		if conversation.FollowUpDate == "" || conversation.Status == "RESOLVED" {
			continue
		}

		followUpDate, err := time.Parse(followUpDateLayout, conversation.FollowUpDate)

		if err != nil {
			continue
		}

		if followUpDate.Before(now) {
			summary.OverdueFollowUps++
		} else if followUpDate.Before(upcomingCutoff) {
			summary.UpcomingFollowUps++
		}
	}

	for _, conversation := range source {
		if conversation.Status == "RESOLVED" {
			continue
		}

		suggested, _ := suggestPriority(conversation, now)

		if priorityRank(suggested) > priorityRank(conversation.Priority) {
			summary.SuggestedPriorityBumps++
		}
	}

	return summary
}
