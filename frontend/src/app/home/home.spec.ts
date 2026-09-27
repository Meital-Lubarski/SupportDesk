import { TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { Home } from './home';
import { DashboardSummary } from '../conversation';
import { ConversationService } from '../conversation.service';

describe('Home', () => {
  const mockSummary: DashboardSummary = {
    total: 4,
    byStatus: {},
    byPriority: {},
    byTag: {},
    byAgent: {},
    overdueFollowUps: 1,
    upcomingFollowUps: 2,
    suggestedPriorityBumps: 1,
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Home],
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

  it('renders a snapshot of the dashboard summary', () => {
    const fixture = TestBed.createComponent(Home);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;

    expect(page.textContent).toContain('4');
    expect(page.textContent).toContain('Total conversations');
    expect(page.textContent).toContain('Overdue follow-ups');
  });

  it('emits navigate when a call-to-action button is clicked', () => {
    const fixture = TestBed.createComponent(Home);
    fixture.detectChanges();

    const emitted: string[] = [];
    fixture.componentInstance.navigate.subscribe((view) => emitted.push(view));

    const page = fixture.nativeElement as HTMLElement;
    const buttons = Array.from(page.querySelectorAll('button'));
    const conversationsButton = buttons.find(
      (b) => b.textContent?.trim() === 'View conversations',
    ) as HTMLButtonElement;

    conversationsButton.click();

    expect(emitted).toEqual(['conversations']);
  });
});
