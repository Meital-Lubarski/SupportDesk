package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// One recorded change on a conversation.
type ActivityEntry struct {
	Timestamp   string `json:"timestamp"`
	Description string `json:"description"`
}

// Keyed by conversation ID. Like `conversations`, this is a plain in-memory
// map with no mutex — fine for a single local demo process, not for
// concurrent traffic.
var activityLog = map[string][]ActivityEntry{}

func init() {
	for _, conversation := range conversations {
		activityLog[conversation.ID] = []ActivityEntry{
			{Timestamp: conversation.CreatedAt, Description: "Conversation created"},
		}
	}

	activityLog["1"] = append(activityLog["1"],
		ActivityEntry{Timestamp: "2026-09-24T10:00:00Z", Description: "Assigned to Alex Chen"},
		ActivityEntry{Timestamp: "2026-09-26T14:15:00Z", Description: "Tags updated to: VIP, Login"},
	)

	activityLog["2"] = append(activityLog["2"],
		ActivityEntry{Timestamp: "2026-09-20T09:00:00Z", Description: "Assigned to Priya Patel"},
	)
}

// Returns the activity log for a single conversation, looked up by ID.
func conversationActivityHandler(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if findConversationIndex(id) == -1 {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	entries := activityLog[id]

	if entries == nil {
		entries = []ActivityEntry{}
	}

	writeJSON(w, http.StatusOK, entries)
}

// Appends one entry per change between before and after.
func recordActivity(id string, before Conversation, after Conversation, now time.Time) {
	for _, description := range diffActivity(before, after) {
		activityLog[id] = append(activityLog[id], ActivityEntry{
			Timestamp:   now.Format(time.RFC3339),
			Description: description,
		})
	}
}

// Compares before/after and returns a human-readable description per change.
func diffActivity(before Conversation, after Conversation) []string {
	changes := make([]string, 0, 5)

	if before.Status != after.Status {
		changes = append(changes, fmt.Sprintf(
			"Status changed from %s to %s", statusPhrase(before.Status), statusPhrase(after.Status),
		))
	}

	if before.Priority != after.Priority {
		changes = append(changes, fmt.Sprintf(
			"Priority changed from %s to %s", strings.ToLower(before.Priority), strings.ToLower(after.Priority),
		))
	}

	if before.AssignedTo != after.AssignedTo {
		changes = append(changes, describeAssigneeChange(before.AssignedTo, after.AssignedTo))
	}

	if before.FollowUpDate != after.FollowUpDate {
		changes = append(changes, describeFollowUpChange(before.FollowUpDate, after.FollowUpDate))
	}

	if !equalTags(before.Tags, after.Tags) {
		changes = append(changes, fmt.Sprintf("Tags updated to: %s", tagsOrNone(after.Tags)))
	}

	return changes
}

func describeAssigneeChange(before string, after string) string {
	switch {
	case before == "":
		return fmt.Sprintf("Assigned to %s", after)
	case after == "":
		return fmt.Sprintf("Unassigned from %s", before)
	default:
		return fmt.Sprintf("Reassigned from %s to %s", before, after)
	}
}

func describeFollowUpChange(before string, after string) string {
	switch {
	case before == "":
		return fmt.Sprintf("Follow-up date set to %s", after)
	case after == "":
		return "Follow-up date cleared"
	default:
		return fmt.Sprintf("Follow-up date changed from %s to %s", before, after)
	}
}

func equalTags(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func tagsOrNone(tags []string) string {
	if len(tags) == 0 {
		return "none"
	}

	return strings.Join(tags, ", ")
}
