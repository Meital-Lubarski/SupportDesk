package main

import (
	"encoding/csv"
	"net/http"
	"strings"
)

// Streams the filtered conversations as a downloadable CSV file. Accepts the
// same query parameters as GET /api/conversations, so an export always
// matches what the UI has filtered down to.
func conversationsExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filtered := filterConversations(r.URL.Query())

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="conversations.csv"`)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"ID", "Customer Name", "Customer Email", "Subject", "Status", "Priority",
		"Assigned To", "Tags", "Follow-up Date", "Notes", "Created At",
	})

	for _, conversation := range filtered {
		writer.Write([]string{
			conversation.ID,
			conversation.CustomerName,
			conversation.CustomerEmail,
			conversation.Subject,
			conversation.Status,
			conversation.Priority,
			conversation.AssignedTo,
			strings.Join(conversation.Tags, "; "),
			conversation.FollowUpDate,
			conversation.Notes,
			conversation.CreatedAt,
		})
	}
}
