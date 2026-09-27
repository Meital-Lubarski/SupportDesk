import { TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { ConversationDetail } from './conversation-detail';
import { ActivityEntry, Conversation, ConversationInsights } from '../conversation';
import { ConversationService } from '../conversation.service';

describe('ConversationDetail', () => {
  const conversation: Conversation = {
    id: '1',
    customerName: 'John Smith',
    customerEmail: 'john.smith@example.com',
    subject: 'Unable to reset password',
    status: 'OPEN',
    priority: 'LOW',
    createdAt: '2026-09-10T09:30:00Z',
    tags: [],
    followUpDate: '',
    notes: '',
    assignedTo: '',
  };

  const insights: ConversationInsights = {
    summary: 'Summary text.',
    suggestedPriority: 'HIGH',
    priorityReason: 'Based on: mentions "urgent".',
    matchesCurrent: false,
  };

  const activityEntries: ActivityEntry[] = [
    { timestamp: '2026-09-10T09:30:00Z', description: 'Conversation created' },
    { timestamp: '2026-09-24T10:00:00Z', description: 'Assigned to Alex Chen' },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConversationDetail],
      providers: [
        {
          provide: ConversationService,
          useValue: {
            getInsights: () => of(insights),
            getActivity: () => of(activityEntries),
          },
        },
      ],
    }).compileComponents();
  });

  it('applies the AI-suggested priority when requested', () => {
    const fixture = TestBed.createComponent(ConversationDetail);
    fixture.componentRef.setInput('conversation', conversation);
    fixture.detectChanges();

    const component = fixture.componentInstance;
    expect(component.editPriority).toBe('LOW');

    component.applySuggestedPriority();

    expect(component.editPriority).toBe('HIGH');
  });

  it('shows the activity history newest first', () => {
    const fixture = TestBed.createComponent(ConversationDetail);
    fixture.componentRef.setInput('conversation', conversation);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    const descriptions = Array.from(page.querySelectorAll('.activity-description')).map((el) =>
      el.textContent?.trim(),
    );

    expect(descriptions).toEqual(['Assigned to Alex Chen', 'Conversation created']);
  });
});
