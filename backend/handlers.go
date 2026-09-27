package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Contains the fields that may be changed by PATCH.
type UpdateConversationRequest struct {
	Status       *string   `json:"status"`
	Priority     *string   `json:"priority"`
	Tags         *[]string `json:"tags"`
	FollowUpDate *string   `json:"followUpDate"`
	Notes        *string   `json:"notes"`
	AssignedTo   *string   `json:"assignedTo"`
}

// Contains the fields accepted when creating a conversation.
type CreateConversationRequest struct {
	CustomerName  string   `json:"customerName"`
	CustomerEmail string   `json:"customerEmail"`
	Subject       string   `json:"subject"`
	Priority      *string  `json:"priority"`
	AssignedTo    *string  `json:"assignedTo"`
	Tags          []string `json:"tags"`
}

const followUpDateLayout = "2006-01-02"

// Returns or creates conversations, depending on the method.
func conversationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, filterConversations(r.URL.Query()))

	case http.MethodPost:
		createConversation(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Validates the request and appends a new conversation, open by default.
func createConversation(w http.ResponseWriter, r *http.Request) {
	var request CreateConversationRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	customerName := strings.TrimSpace(request.CustomerName)
	customerEmail := strings.TrimSpace(request.CustomerEmail)
	subject := strings.TrimSpace(request.Subject)

	if customerName == "" || customerEmail == "" || subject == "" {
		http.Error(w, "customerName, customerEmail, and subject are required", http.StatusBadRequest)
		return
	}

	if !strings.Contains(customerEmail, "@") {
		http.Error(w, "customerEmail must be a valid email address", http.StatusBadRequest)
		return
	}

	priority := "MEDIUM"

	if request.Priority != nil {
		trimmed := strings.ToUpper(strings.TrimSpace(*request.Priority))

		if trimmed != "" {
			if !isValidPriority(trimmed) {
				http.Error(w, "Invalid priority", http.StatusBadRequest)
				return
			}

			priority = trimmed
		}
	}

	assignedTo := ""

	if request.AssignedTo != nil {
		trimmed := strings.TrimSpace(*request.AssignedTo)

		if !isValidAssignee(trimmed) {
			http.Error(w, "Invalid assignedTo", http.StatusBadRequest)
			return
		}

		assignedTo = trimmed
	}

	newConversation := Conversation{
		ID:            nextConversationID(),
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		Subject:       subject,
		Status:        "OPEN",
		Priority:      priority,
		CreatedAt:     time.Now().Format(time.RFC3339),
		Tags:          normalizeTags(request.Tags),
		AssignedTo:    assignedTo,
	}

	conversations = append(conversations, newConversation)
	activityLog[newConversation.ID] = []ActivityEntry{
		{Timestamp: newConversation.CreatedAt, Description: "Conversation created"},
	}

	writeJSON(w, http.StatusCreated, newConversation)
}

// Returns the next sequential ID, based on the highest existing numeric ID.
func nextConversationID() string {
	maxID := 0

	for _, conversation := range conversations {
		if id, err := strconv.Atoi(conversation.ID); err == nil && id > maxID {
			maxID = id
		}
	}

	return strconv.Itoa(maxID + 1)
}

// Shared by the JSON list endpoint and the CSV export endpoint, so both
// always agree on what a given set of filters matches.
func filterConversations(query url.Values) []Conversation {
	statusFilter := strings.ToUpper(strings.TrimSpace(query.Get("status")))
	priorityFilter := strings.ToUpper(strings.TrimSpace(query.Get("priority")))
	searchTerm := strings.ToLower(strings.TrimSpace(query.Get("search")))
	tagFilter := strings.ToLower(strings.TrimSpace(query.Get("tag")))
	assigneeFilter := strings.TrimSpace(query.Get("assignee"))

	filtered := make([]Conversation, 0)

	for _, conversation := range conversations {
		if statusFilter != "" && conversation.Status != statusFilter {
			continue
		}

		if priorityFilter != "" && conversation.Priority != priorityFilter {
			continue
		}

		if tagFilter != "" && !hasTag(conversation.Tags, tagFilter) {
			continue
		}

		if assigneeFilter != "" && !matchesAssignee(conversation.AssignedTo, assigneeFilter) {
			continue
		}

		if searchTerm != "" {
			name := strings.ToLower(conversation.CustomerName)
			email := strings.ToLower(conversation.CustomerEmail)
			subject := strings.ToLower(conversation.Subject)

			matchesSearch := strings.Contains(name, searchTerm) ||
				strings.Contains(email, searchTerm) ||
				strings.Contains(subject, searchTerm)

			if !matchesSearch {
				continue
			}
		}

		filtered = append(filtered, conversation)
	}

	return filtered
}

// Updates a single conversation by ID.
func conversationByIDHandler(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/conversations/")

	if insightsID, ok := strings.CutSuffix(rest, "/insights"); ok {
		conversationInsightsHandler(w, r, insightsID)
		return
	}

	if activityID, ok := strings.CutSuffix(rest, "/activity"); ok {
		conversationActivityHandler(w, r, activityID)
		return
	}

	id := rest

	if id == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	index := findConversationIndex(id)

	if index == -1 {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, conversations[index])

	case http.MethodPatch:
		updateConversation(w, r, index)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func findConversationIndex(id string) int {
	for index, conversation := range conversations {
		if conversation.ID == id {
			return index
		}
	}

	return -1
}

// Validates changes on a copy before updating the stored conversation.
func updateConversation(w http.ResponseWriter, r *http.Request, index int) {
	var update UpdateConversationRequest

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	before := conversations[index]
	updatedConversation := before

	if update.Status != nil {
		status := strings.ToUpper(strings.TrimSpace(*update.Status))

		if !isValidStatus(status) {
			http.Error(w, "Invalid status; no changes applied", http.StatusBadRequest)
			return
		}

		updatedConversation.Status = status
	}

	if update.Priority != nil {
		priority := strings.ToUpper(strings.TrimSpace(*update.Priority))

		if !isValidPriority(priority) {
			http.Error(w, "Invalid priority; no changes applied", http.StatusBadRequest)
			return
		}

		updatedConversation.Priority = priority
	}

	if update.FollowUpDate != nil {
		followUpDate := strings.TrimSpace(*update.FollowUpDate)

		if followUpDate != "" {
			if _, err := time.Parse(followUpDateLayout, followUpDate); err != nil {
				http.Error(w, "Invalid followUpDate; expected YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}

		updatedConversation.FollowUpDate = followUpDate
	}

	if update.Notes != nil {
		updatedConversation.Notes = *update.Notes
	}

	if update.Tags != nil {
		updatedConversation.Tags = normalizeTags(*update.Tags)
	}

	if update.AssignedTo != nil {
		assignedTo := strings.TrimSpace(*update.AssignedTo)

		if !isValidAssignee(assignedTo) {
			http.Error(w, "Invalid assignedTo; no changes applied", http.StatusBadRequest)
			return
		}

		updatedConversation.AssignedTo = assignedTo
	}

	conversations[index] = updatedConversation
	recordActivity(updatedConversation.ID, before, updatedConversation, time.Now())
	writeJSON(w, http.StatusOK, updatedConversation)
}

// Returns true if assignedTo matches the filter. The special value
// "unassigned" (case-insensitive) matches conversations with no assignee.
func matchesAssignee(assignedTo string, filter string) bool {
	if strings.EqualFold(filter, "unassigned") {
		return assignedTo == ""
	}

	return strings.EqualFold(assignedTo, filter)
}

// Returns true if tags contains tag, case-insensitively.
func hasTag(tags []string, tag string) bool {
	for _, candidate := range tags {
		if strings.ToLower(candidate) == tag {
			return true
		}
	}

	return false
}

// Trims, drops empty entries, and de-duplicates tags case-insensitively.
func normalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	normalized := make([]string, 0, len(tags))

	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)

		if trimmed == "" {
			continue
		}

		key := strings.ToLower(trimmed)

		if seen[key] {
			continue
		}

		seen[key] = true
		normalized = append(normalized, trimmed)
	}

	return normalized
}

func isValidStatus(status string) bool {
	return status == "OPEN" ||
		status == "IN_PROGRESS" ||
		status == "RESOLVED"
}

func isValidPriority(priority string) bool {
	return priority == "LOW" ||
		priority == "MEDIUM" ||
		priority == "HIGH"
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
