//A single conversation
export interface Conversation {
  id: string;
  customerName: string;
  customerEmail: string;
  subject: string;
  status: string;
  priority: string;
  createdAt: string;
  tags: string[];
  followUpDate: string;
  notes: string;
  assignedTo: string;
}

//The data we send when we update the conversation
export interface UpdateConversationRequest {
  status?: string;
  priority?: string;
  tags?: string[];
  followUpDate?: string;
  notes?: string;
  assignedTo?: string;
}

//The data we send when we create a conversation
export interface CreateConversationRequest {
  customerName: string;
  customerEmail: string;
  subject: string;
  priority?: string;
  assignedTo?: string;
  tags?: string[];
}

//Aggregated counts served by the dashboard endpoint
export interface DashboardSummary {
  total: number;
  byStatus: Record<string, number>;
  byPriority: Record<string, number>;
  byTag: Record<string, number>;
  byAgent: Record<string, number>;
  overdueFollowUps: number;
  upcomingFollowUps: number;
  suggestedPriorityBumps: number;
}

//Rule-based "smart" read on a single conversation (see backend/insights.go)
export interface ConversationInsights {
  summary: string;
  suggestedPriority: string;
  priorityReason: string;
  matchesCurrent: boolean;
}

//One recorded change on a conversation (see backend/activity.go)
export interface ActivityEntry {
  timestamp: string;
  description: string;
}