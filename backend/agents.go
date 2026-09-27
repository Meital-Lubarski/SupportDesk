package main

import "net/http"

// The fixed support-agent roster. A real system would back this with a user
// directory; this project keeps everything in-memory, so a fixed list matches
// the rest of the app's storage model.
var agents = []string{"Alex Chen", "Priya Patel", "Jordan Lee", "Sam Rivera"}

// Returns true if name is empty (unassigned) or one of the known agents.
func isValidAssignee(name string) bool {
	if name == "" {
		return true
	}

	for _, agent := range agents {
		if agent == name {
			return true
		}
	}

	return false
}

// Returns the agent roster.
func agentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, agents)
}
