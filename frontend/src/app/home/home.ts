import { Component, EventEmitter, OnInit, Output, signal } from '@angular/core';
import { DashboardSummary } from '../conversation';
import { ConversationService } from '../conversation.service';
import { AppView } from '../shared/app-view';

interface Highlight {
  icon: string;
  title: string;
  description: string;
}

//Landing screen: a quick snapshot plus entry points into the rest of the app.
@Component({
  selector: 'app-home',
  imports: [],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit {
  @Output() navigate = new EventEmitter<AppView>();

  summary = signal<DashboardSummary | null>(null);

  readonly highlights: Highlight[] = [
    {
      icon: '🏷️',
      title: 'Tags & assignees',
      description: 'Categorize conversations and route them to the right agent.',
    },
    {
      icon: '📅',
      title: 'Follow-ups',
      description: 'Track follow-up dates with automatic overdue and due-soon flags.',
    },
    {
      icon: '✨',
      title: 'AI Insights',
      description: 'A rule-based summary and suggested priority for every conversation.',
    },
    {
      icon: '📊',
      title: 'Live dashboard',
      description: 'Status, priority, tag, and workload breakdowns at a glance.',
    },
  ];

  constructor(private conversationService: ConversationService) {}

  ngOnInit(): void {
    this.conversationService.getDashboardSummary().subscribe({
      next: (summary) => this.summary.set(summary),
      error: (error) => console.error('Failed to load dashboard snapshot:', error),
    });
  }

  goTo(view: AppView): void {
    this.navigate.emit(view);
  }
}
