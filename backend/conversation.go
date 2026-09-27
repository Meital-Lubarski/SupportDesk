package main

// A customer conversation.
type Conversation struct {
	ID            string   `json:"id"`
	CustomerName  string   `json:"customerName"`
	CustomerEmail string   `json:"customerEmail"`
	Subject       string   `json:"subject"`
	Status        string   `json:"status"`
	Priority      string   `json:"priority"`
	CreatedAt     string   `json:"createdAt"`
	Tags          []string `json:"tags"`
	FollowUpDate  string   `json:"followUpDate"`
	Notes         string   `json:"notes"`
	AssignedTo    string   `json:"assignedTo"`
}

// The in-memory store.
var conversations = []Conversation{
	{
		ID:            "1",
		CustomerName:  "John Smith",
		CustomerEmail: "john.smith@example.com",
		Subject:       "Unable to reset password",
		Status:        "OPEN",
		Priority:      "HIGH",
		CreatedAt:     "2026-09-10T09:30:00Z",
		Tags:          []string{"VIP", "Login"},
		FollowUpDate:  "2026-09-29",
		Notes:         "Customer locked out on mobile app; sent temporary reset link.",
		AssignedTo:    "Alex Chen",
	},
	{
		ID:            "2",
		CustomerName:  "Sarah Cohen",
		CustomerEmail: "sarah.cohen@example.com",
		Subject:       "Question about billing",
		Status:        "IN_PROGRESS",
		Priority:      "MEDIUM",
		CreatedAt:     "2026-09-11T12:15:00Z",
		Tags:          []string{"Billing"},
		FollowUpDate:  "2026-09-24",
		Notes:         "Waiting on invoice from finance team before responding.",
		AssignedTo:    "Priya Patel",
	},
	{
		ID:            "3",
		CustomerName:  "Daniel Levi",
		CustomerEmail: "daniel.levi@example.com",
		Subject:       "Account access restored",
		Status:        "RESOLVED",
		Priority:      "LOW",
		CreatedAt:     "2026-09-12T15:45:00Z",
		Tags:          []string{"Account"},
		FollowUpDate:  "",
		Notes:         "",
		AssignedTo:    "Priya Patel",
	},
	{
		ID:            "4",
		CustomerName:  "Emily Johnson",
		CustomerEmail: "emily.johnson@example.com",
		Subject:       "Application crashes on login",
		Status:        "OPEN",
		Priority:      "HIGH",
		CreatedAt:     "2026-09-13T08:20:00Z",
		Tags:          []string{"Bug", "VIP"},
		FollowUpDate:  "2026-09-28",
		Notes:         "Reproduced on Android 14; escalated to mobile team.",
	},
}
