import { Component, OnInit, signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import {
  Conversation,
  CreateConversationRequest,
  UpdateConversationRequest,
} from '../conversation';
import { ConversationFilters, ConversationService } from '../conversation.service';
import { ConversationDetail } from '../conversation-detail/conversation-detail';
import { NewConversation } from '../new-conversation/new-conversation';
import {
  priorityBadgeClass,
  priorityLabel,
  statusBadgeClass,
  statusLabel,
} from '../shared/badge.util';
import { followUpState } from '../shared/follow-up.util';

//Handles browsing, filtering, and editing conversations.
@Component({
  selector: 'app-conversations',
  imports: [DatePipe, ConversationDetail, NewConversation],
  templateUrl: './conversations.html',
  styleUrl: './conversations.css',
})
export class Conversations implements OnInit {
  conversations = signal<Conversation[]>([]);
  selectedConversation = signal<Conversation | null>(null);

  searchTerm = signal('');
  statusFilter = signal('');
  priorityFilter = signal('');
  tagFilter = signal('');
  assigneeFilter = signal('');

  agents = signal<string[]>([]);

  loading = signal(true);
  errorMessage = signal('');

  saving = signal(false);
  saveMessage = signal('');
  saveError = signal('');

  showCreateForm = signal(false);
  creating = signal(false);
  createError = signal('');

  readonly statusBadgeClass = statusBadgeClass;
  readonly priorityBadgeClass = priorityBadgeClass;
  readonly statusLabel = statusLabel;
  readonly priorityLabel = priorityLabel;
  readonly followUpState = followUpState;

  constructor(private conversationService: ConversationService) {}

  ngOnInit(): void {
    this.loadConversations();

    this.conversationService.getAgents().subscribe({
      next: (agents) => this.agents.set(agents),
      error: (error) => console.error('Failed to load agents:', error),
    });
  }

  //The filters currently applied to the list, as loaded from the backend.
  private activeFilters(): ConversationFilters {
    return {
      status: this.statusFilter(),
      priority: this.priorityFilter(),
      search: this.searchTerm(),
      tag: this.tagFilter(),
      assignee: this.assigneeFilter(),
    };
  }

  //A download link for the currently applied filters, so the export always
  //matches what's on screen.
  exportUrl(): string {
    return this.conversationService.getExportUrl(this.activeFilters());
  }

  //Loads the conversations using the active filters.
  loadConversations(): void {
    this.loading.set(true);
    this.errorMessage.set('');

    this.conversationService
      .getConversations(this.activeFilters())
      .subscribe({
        next: (conversations) => {
          this.conversations.set(conversations);
          this.loading.set(false);

          const selected = this.selectedConversation();

          if (selected) {
            const refreshed = conversations.find((c) => c.id === selected.id);
            this.selectedConversation.set(refreshed ?? null);
          }
        },
        error: (error) => {
          console.error('Failed to load conversations:', error);
          this.errorMessage.set('Failed to load conversations.');
          this.loading.set(false);
        },
      });
  }

  //Expands the clicked conversation's edit form, or collapses it if already open.
  toggleConversation(conversation: Conversation): void {
    const isAlreadyOpen = this.selectedConversation()?.id === conversation.id;
    this.selectedConversation.set(isAlreadyOpen ? null : conversation);
    this.saveMessage.set('');
    this.saveError.set('');
  }

  applyFilters(
    search: string,
    status: string,
    priority: string,
    tag: string,
    assignee: string,
  ): void {
    this.searchTerm.set(search);
    this.statusFilter.set(status);
    this.priorityFilter.set(priority);
    this.tagFilter.set(tag);
    this.assigneeFilter.set(assignee);
    this.selectedConversation.set(null);
    this.saveMessage.set('');
    this.saveError.set('');

    this.loadConversations();
  }

  //Resets every filter back to its default.
  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('');
    this.priorityFilter.set('');
    this.tagFilter.set('');
    this.assigneeFilter.set('');
    this.selectedConversation.set(null);
    this.saveMessage.set('');
    this.saveError.set('');

    this.loadConversations();
  }

  //Turns "Alex Chen" into "AC" for the avatar chip.
  initials(name: string): string {
    return name
      .split(' ')
      .filter(Boolean)
      .map((part) => part[0])
      .join('')
      .slice(0, 2)
      .toUpperCase();
  }

  openCreateForm(): void {
    this.createError.set('');
    this.showCreateForm.set(true);
  }

  closeCreateForm(): void {
    this.showCreateForm.set(false);
  }

  //Creates a conversation, then selects it in the (now refreshed) list.
  createConversation(request: CreateConversationRequest): void {
    this.creating.set(true);
    this.createError.set('');

    this.conversationService.createConversation(request).subscribe({
      next: (created) => {
        this.creating.set(false);
        this.showCreateForm.set(false);
        this.selectedConversation.set(created);
        this.loadConversations();
      },
      error: (error) => {
        console.error('Failed to create conversation:', error);
        this.creating.set(false);
        this.createError.set('Failed to create conversation. Check that all fields are valid.');
      },
    });
  }

  //Saves edits from the detail panel to the backend.
  saveConversation(update: UpdateConversationRequest): void {
    const selected = this.selectedConversation();

    if (!selected) {
      return;
    }

    this.saving.set(true);
    this.saveMessage.set('');
    this.saveError.set('');

    this.conversationService
      .updateConversation(selected.id, update)
      .subscribe({
        next: (updatedConversation) => {
          this.selectedConversation.set(updatedConversation);
          this.saving.set(false);
          this.saveMessage.set('Changes saved successfully.');
          this.loadConversations();
        },
        error: (error) => {
          console.error('Failed to update conversation:', error);
          this.saving.set(false);
          this.saveError.set('Failed to save changes.');
        },
      });
  }
}
