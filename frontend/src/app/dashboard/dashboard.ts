import { Component, OnInit, signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import { DashboardSummary } from '../conversation';
import { ConversationService } from '../conversation.service';
import { priorityLabel, statusLabel } from '../shared/badge.util';

interface CountEntry {
  key: string;
  label: string;
  count: number;
  percent: number;
}

//Shows aggregated conversation counts.
@Component({
  selector: 'app-dashboard',
  imports: [DatePipe],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  summary = signal<DashboardSummary | null>(null);
  loading = signal(true);
  errorMessage = signal('');
  lastUpdated = signal<Date | null>(null);

  //Fixed workflow order — not alphabetical, so the chart always reads left to right as a pipeline.
  private readonly statusOrder = ['OPEN', 'IN_PROGRESS', 'RESOLVED'];

  //Fixed severity order — Low to High, never sorted by count.
  private readonly priorityOrder = ['LOW', 'MEDIUM', 'HIGH'];

  //Reuses the same status/priority colors as the badges elsewhere in the app.
  readonly statusColors: Record<string, string> = {
    OPEN: 'var(--color-status-open)',
    IN_PROGRESS: 'var(--color-status-in-progress)',
    RESOLVED: 'var(--color-status-resolved)',
  };

  constructor(private conversationService: ConversationService) {}

  ngOnInit(): void {
    this.loadSummary();
  }

  loadSummary(): void {
    this.loading.set(true);
    this.errorMessage.set('');

    this.conversationService.getDashboardSummary().subscribe({
      next: (summary) => {
        this.summary.set(summary);
        this.loading.set(false);
        this.lastUpdated.set(new Date());
      },
      error: (error) => {
        console.error('Failed to load dashboard summary:', error);
        this.errorMessage.set('Failed to load dashboard data.');
        this.loading.set(false);
      },
    });
  }

  //Ordered for the stacked-bar chart (workflow order, zero-count statuses dropped).
  statusChart(): CountEntry[] {
    return this.toOrderedEntries(this.summary()?.byStatus, this.statusOrder, statusLabel);
  }

  priorityBreakdown(): CountEntry[] {
    return this.toOrderedEntries(this.summary()?.byPriority, this.priorityOrder, priorityLabel);
  }

  tagBreakdown(): CountEntry[] {
    return this.toEntries(this.summary()?.byTag, (tag) => tag).sort(
      (a, b) => b.count - a.count,
    );
  }

  //Busiest agents first; "Unassigned" always trails, even if it has a high count.
  agentBreakdown(): CountEntry[] {
    return this.toEntries(this.summary()?.byAgent, (agent) => agent).sort((a, b) => {
      if (a.key === 'Unassigned') return 1;
      if (b.key === 'Unassigned') return -1;
      return b.count - a.count;
    });
  }

  private toOrderedEntries(
    counts: Record<string, number> | undefined,
    order: string[],
    labelFor: (key: string) => string,
  ): CountEntry[] {
    if (!counts) {
      return [];
    }

    const total = this.summary()?.total || 1;

    return order
      .filter((key) => (counts[key] ?? 0) > 0)
      .map((key) => ({
        key,
        label: labelFor(key),
        count: counts[key],
        percent: Math.round((counts[key] / total) * 100),
      }));
  }

  private toEntries(
    counts: Record<string, number> | undefined,
    labelFor: (key: string) => string,
  ): CountEntry[] {
    if (!counts) {
      return [];
    }

    const total = this.summary()?.total || 1;

    return Object.entries(counts).map(([key, count]) => ({
      key,
      label: labelFor(key),
      count,
      percent: Math.round((count / total) * 100),
    }));
  }
}
