package summary

// NO-AI TASK
type Conversation struct {
	ID       string
	Customer string
	Priority string
	Status   string
}

func SummarizeConversations(conversations []Conversation) map[string]int {
	summary := make(map[string]int)
	for _, conversation := range conversations {
		if conversation.Status == "" || conversation.Priority == "" {
			continue
		}
		key := conversation.Status + "_" + conversation.Priority
		summary[key]++
	}
	return summary
}
