package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// A rule-based "smart" read on a single conversation. Not backed by a language
// model — it is a deterministic keyword/heuristic engine over local data.
type ConversationInsights struct {
	Summary           string `json:"summary"`
	SuggestedPriority string `json:"suggestedPriority"`
	PriorityReason    string `json:"priorityReason"`
	MatchesCurrent    bool   `json:"matchesCurrent"`
}

var urgentKeywords = []string{
	"crash", "crashes", "crashed", "down", "broken", "cannot", "can't",
	"unable", "urgent", "angry", "refund", "locked out", "not working",
	"failure", "lost access", "security", "breach",
}

var lowUrgencyKeywords = []string{
	"question", "wondering", "curious", "how do i", "just asking",
}

// Returns insights for a single conversation, looked up by ID.
func conversationInsightsHandler(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	index := findConversationIndex(id)

	if index == -1 {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, buildInsights(conversations[index], time.Now()))
}

func buildInsights(c Conversation, now time.Time) ConversationInsights {
	priority, reason := suggestPriority(c, now)

	return ConversationInsights{
		Summary:           summarize(c, now),
		SuggestedPriority: priority,
		PriorityReason:    reason,
		MatchesCurrent:    priority == c.Priority,
	}
}

// Scores urgency signals in the conversation and maps the score to a priority.
func suggestPriority(c Conversation, now time.Time) (string, string) {
	text := strings.ToLower(c.Subject + " " + c.Notes)
	score := 0
	reasons := make([]string, 0, 4)

	for _, keyword := range urgentKeywords {
		if strings.Contains(text, keyword) {
			score += 2
			reasons = append(reasons, fmt.Sprintf("mentions %q", keyword))
			break
		}
	}

	for _, keyword := range lowUrgencyKeywords {
		if strings.Contains(text, keyword) {
			score--
			reasons = append(reasons, "reads like a general question")
			break
		}
	}

	if hasTag(c.Tags, "vip") {
		score += 2
		reasons = append(reasons, "VIP customer")
	}

	if c.Status != "RESOLVED" && c.FollowUpDate != "" {
		if followUpDate, err := time.Parse(followUpDateLayout, c.FollowUpDate); err == nil && followUpDate.Before(now) {
			score++
			reasons = append(reasons, "follow-up is overdue")
		}
	}

	if c.Status == "OPEN" {
		if createdAt, err := time.Parse(time.RFC3339, c.CreatedAt); err == nil {
			daysOpen := int(now.Sub(createdAt).Hours() / 24)

			if daysOpen >= 3 {
				score++
				reasons = append(reasons, fmt.Sprintf("open for %d days", daysOpen))
			}
		}
	}

	priority := "LOW"

	switch {
	case score >= 3:
		priority = "HIGH"
	case score >= 1:
		priority = "MEDIUM"
	}

	reasonText := "No strong urgency signals found."

	if len(reasons) > 0 {
		reasonText = "Based on: " + strings.Join(reasons, ", ") + "."
	}

	return priority, reasonText
}

// Builds a short templated recap of the conversation.
func summarize(c Conversation, now time.Time) string {
	intro := fmt.Sprintf(
		"%s contacted support about \"%s\". Currently %s with %s priority",
		c.CustomerName, c.Subject, statusPhrase(c.Status), strings.ToLower(c.Priority),
	)

	if len(c.Tags) > 0 {
		intro += fmt.Sprintf(", tagged %s", strings.Join(c.Tags, ", "))
	}

	intro += "."

	sentences := []string{intro}

	if followUp := followUpPhrase(c, now); followUp != "" {
		sentences = append(sentences, followUp)
	}

	if c.Notes != "" {
		sentences = append(sentences, fmt.Sprintf("Latest notes: %q.", c.Notes))
	} else {
		sentences = append(sentences, "No internal notes yet.")
	}

	return strings.Join(sentences, " ")
}

func statusPhrase(status string) string {
	switch status {
	case "OPEN":
		return "open"
	case "IN_PROGRESS":
		return "in progress"
	case "RESOLVED":
		return "resolved"
	default:
		return strings.ToLower(status)
	}
}

func followUpPhrase(c Conversation, now time.Time) string {
	if c.FollowUpDate == "" || c.Status == "RESOLVED" {
		return ""
	}

	followUpDate, err := time.Parse(followUpDateLayout, c.FollowUpDate)

	if err != nil {
		return ""
	}

	switch {
	case followUpDate.Before(now):
		return fmt.Sprintf("A follow-up was due on %s and is now overdue.", c.FollowUpDate)
	case followUpDate.Before(now.Add(upcomingFollowUpWindow)):
		return fmt.Sprintf("A follow-up is scheduled for %s.", c.FollowUpDate)
	default:
		return fmt.Sprintf("Next follow-up is set for %s.", c.FollowUpDate)
	}
}

// Returns a numeric rank so priorities can be compared.
func priorityRank(priority string) int {
	switch priority {
	case "HIGH":
		return 2
	case "MEDIUM":
		return 1
	default:
		return 0
	}
}
