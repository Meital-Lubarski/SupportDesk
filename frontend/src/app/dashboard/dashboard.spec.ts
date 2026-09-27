import { TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { Dashboard } from './dashboard';
import { DashboardSummary } from '../conversation';
import { ConversationService } from '../conversation.service';

describe('Dashboard', () => {
  const mockSummary: DashboardSummary = {
    total: 4,
    byStatus: { OPEN: 2, IN_PROGRESS: 1, RESOLVED: 1 },
    byPriority: { HIGH: 2, MEDIUM: 1, LOW: 1 },
    byTag: { VIP: 2, Billing: 1 },
    byAgent: { 'Alex Chen': 2, Unassigned: 1 },
    overdueFollowUps: 1,
    upcomingFollowUps: 2,
    suggestedPriorityBumps: 1,
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Dashboard],
      providers: [
        {
          provide: ConversationService,
          useValue: {
            getDashboardSummary: () => of(mockSummary),
          },
        },
      ],
    }).compileComponents();
  });

  it('renders aggregated counts from the service', () => {
    const fixture = TestBed.createComponent(Dashboard);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;

    expect(page.textContent).toContain('4');
    expect(page.textContent).toContain('Total conversations');
    expect(page.textContent).toContain('Overdue follow-ups');
    expect(page.textContent).toContain('AI-flagged priority bumps');
    expect(page.textContent).toContain('VIP');
  });

  it('renders the status distribution chart in workflow order with a legend', () => {
    const fixture = TestBed.createComponent(Dashboard);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;

    const segments = Array.from(page.querySelectorAll('.stack-segment'));
    expect(segments.length).toBe(3);

    const legendLabels = Array.from(page.querySelectorAll('.legend-label')).map(
      (el) => el.textContent?.trim(),
    );
    expect(legendLabels).toEqual(['Open', 'In progress', 'Resolved']);

    const firstSegment = segments[0] as HTMLElement;
    expect(firstSegment.title).toContain('Open: 2 (50%)');
  });

  it('renders agent workload with Unassigned trailing regardless of count', () => {
    const fixture = TestBed.createComponent(Dashboard);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    const agentEntries = Array.from(page.querySelectorAll('.breakdown-card')).find((card) =>
      card.querySelector('h2')?.textContent?.includes('Agent workload'),
    );

    expect(agentEntries).toBeTruthy();

    const rowLabels = Array.from(
      agentEntries!.querySelectorAll('.breakdown-row span:first-child'),
    ).map((el) => el.textContent?.trim());

    expect(rowLabels).toEqual(['Alex Chen', 'Unassigned']);
  });
});
