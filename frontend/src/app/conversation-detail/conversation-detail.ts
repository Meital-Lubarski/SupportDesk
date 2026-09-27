import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';
import {
  ActivityEntry,
  Conversation,
  ConversationInsights,
  UpdateConversationRequest,
} from '../conversation';
import { ConversationService } from '../conversation.service';
import { priorityBadgeClass, priorityLabel } from '../shared/badge.util';
import { TagEditor } from '../shared/tag-editor/tag-editor';

//Editable panel for the selected conversation.
@Component({
  selector: 'app-conversation-detail',
  imports: [FormsModule, DatePipe, TagEditor],
  templateUrl: './conversation-detail.html',
  styleUrl: './conversation-detail.css',
})
export class ConversationDetail implements OnChanges {
  @Input({ required: true }) conversation!: Conversation;
  @Input() agents: string[] = [];
  @Input() saving = false;
  @Input() saveMessage = '';
  @Input() saveError = '';

  @Output() save = new EventEmitter<UpdateConversationRequest>();

  insights = signal<ConversationInsights | null>(null);
  insightsLoading = signal(true);
  insightsError = signal('');

  activity = signal<ActivityEntry[]>([]);
  activityLoading = signal(true);
  activityError = signal('');

  readonly priorityBadgeClass = priorityBadgeClass;
  readonly priorityLabel = priorityLabel;

  editStatus = '';
  editPriority = '';
  editFollowUpDate = '';
  editNotes = '';
  editTags: string[] = [];
  editAssignedTo = '';

  constructor(private conversationService: ConversationService) {}

  //Angular calls this whenever the @Input conversation reference changes.
  ngOnChanges(): void {
    this.editStatus = this.conversation.status;
    this.editPriority = this.conversation.priority;
    this.editFollowUpDate = this.conversation.followUpDate;
    this.editNotes = this.conversation.notes;
    this.editTags = [...this.conversation.tags];
    this.editAssignedTo = this.conversation.assignedTo;

    this.loadInsights();
    this.loadActivity();
  }

  loadActivity(): void {
    this.activityLoading.set(true);
    this.activityError.set('');

    this.conversationService.getActivity(this.conversation.id).subscribe({
      next: (entries) => {
        //Newest first, so the most recent change is always visible without scrolling.
        this.activity.set([...entries].reverse());
        this.activityLoading.set(false);
      },
      error: (error) => {
        console.error('Failed to load activity:', error);
        this.activityError.set('Activity history is unavailable right now.');
        this.activityLoading.set(false);
      },
    });
  }

  loadInsights(): void {
    this.insightsLoading.set(true);
    this.insightsError.set('');

    this.conversationService.getInsights(this.conversation.id).subscribe({
      next: (insights) => {
        this.insights.set(insights);
        this.insightsLoading.set(false);
      },
      error: (error) => {
        console.error('Failed to load AI insights:', error);
        this.insightsError.set('AI insights are unavailable right now.');
        this.insightsLoading.set(false);
      },
    });
  }

  applySuggestedPriority(): void {
    const suggested = this.insights()?.suggestedPriority;

    if (suggested) {
      this.editPriority = suggested;
    }
  }

  submit(): void {
    this.save.emit({
      status: this.editStatus,
      priority: this.editPriority,
      followUpDate: this.editFollowUpDate,
      notes: this.editNotes,
      tags: this.editTags,
      assignedTo: this.editAssignedTo,
    });
  }
}
