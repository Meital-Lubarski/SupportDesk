import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
  ActivityEntry,
  Conversation,
  ConversationInsights,
  CreateConversationRequest,
  DashboardSummary,
  UpdateConversationRequest,
} from './conversation';

export interface ConversationFilters {
  status?: string;
  priority?: string;
  search?: string;
  tag?: string;
  assignee?: string;
}

//service that talks to the backend API.
@Injectable({
  providedIn: 'root',
})
export class ConversationService {
  private readonly baseUrl = 'http://localhost:8081/api';
  private readonly apiUrl = `${this.baseUrl}/conversations`;

  constructor(private http: HttpClient) {}

  // Gets conversations from the server and filters them.
  getConversations(filters: ConversationFilters): Observable<Conversation[]> {
    return this.http.get<Conversation[]>(this.apiUrl, { params: this.buildParams(filters) });
  }

  //Builds the download URL for the currently active filters. A plain link/new-tab
  //navigation to this URL triggers the browser's normal file download, no fetch needed.
  getExportUrl(filters: ConversationFilters): string {
    const query = this.buildParams(filters).toString();
    return query ? `${this.apiUrl}/export?${query}` : `${this.apiUrl}/export`;
  }

  private buildParams(filters: ConversationFilters): HttpParams {
    let params = new HttpParams();

    if (filters.status) {
      params = params.set('status', filters.status);
    }

    if (filters.priority) {
      params = params.set('priority', filters.priority);
    }

    if (filters.search) {
      params = params.set('search', filters.search);
    }

    if (filters.tag) {
      params = params.set('tag', filters.tag);
    }

    if (filters.assignee) {
      params = params.set('assignee', filters.assignee);
    }

    return params;
  }

  //Gets the fixed list of support agents that conversations can be assigned to.
  getAgents(): Observable<string[]> {
    return this.http.get<string[]>(`${this.baseUrl}/agents`);
  }

  //Sends a POST request to create a new conversation.
  createConversation(request: CreateConversationRequest): Observable<Conversation> {
    return this.http.post<Conversation>(this.apiUrl, request);
  }

  //Sends a PATCH request.
  updateConversation(
    id: string,
    update: UpdateConversationRequest,
  ): Observable<Conversation> {
    return this.http.patch<Conversation>(
      `${this.apiUrl}/${id}`,
      update,
    );
  }

  //Gets aggregated counts for the dashboard.
  getDashboardSummary(): Observable<DashboardSummary> {
    return this.http.get<DashboardSummary>(`${this.baseUrl}/dashboard`);
  }

  //Gets the rule-based summary and suggested priority for one conversation.
  getInsights(id: string): Observable<ConversationInsights> {
    return this.http.get<ConversationInsights>(`${this.apiUrl}/${id}/insights`);
  }

  //Gets the change history for one conversation, oldest first.
  getActivity(id: string): Observable<ActivityEntry[]> {
    return this.http.get<ActivityEntry[]>(`${this.apiUrl}/${id}/activity`);
  }
}