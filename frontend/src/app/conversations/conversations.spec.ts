import { TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';
import { Conversations } from './conversations';
import { Conversation, ConversationInsights } from '../conversation';
import { ConversationService } from '../conversation.service';

describe('Conversations', () => {
  //Fake data to simulate the backend response.
  const mockConversations: Conversation[] = [
    {
      id: '1',
      customerName: 'John Smith',
      customerEmail: 'john.smith@example.com',
      subject: 'Unable to reset password',
      status: 'OPEN',
      priority: 'HIGH',
      createdAt: '2026-09-10T09:30:00Z',
      tags: ['VIP', 'Login'],
      followUpDate: '2026-09-29',
      notes: 'Sent a temporary reset link.',
      assignedTo: 'Alex Chen',
    },
  ];

  const mockInsights: ConversationInsights = {
    summary: 'John Smith contacted support about "Unable to reset password".',
    suggestedPriority: 'HIGH',
    priorityReason: 'Based on: mentions "locked out".',
    matchesCurrent: true,
  };

  const createdConversation: Conversation = {
    id: '2',
    customerName: 'Jane Doe',
    customerEmail: 'jane@example.com',
    subject: 'Cannot log in',
    status: 'OPEN',
    priority: 'MEDIUM',
    createdAt: '2026-09-27T10:00:00Z',
    tags: [],
    followUpDate: '',
    notes: '',
    assignedTo: '',
  };

  let createConversationSpy: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    createConversationSpy = vi.fn().mockReturnValue(of(createdConversation));

    await TestBed.configureTestingModule({
      imports: [Conversations],
      providers: [
        {
          provide: ConversationService,
          useValue: {
            //Includes createdConversation too, mirroring how a real refetch
            //after POST would include the newly created record.
            getConversations: () => of([...mockConversations, createdConversation]),
            getInsights: () => of(mockInsights),
            getActivity: () => of([]),
            getAgents: () => of(['Alex Chen', 'Priya Patel']),
            getExportUrl: (filters: Record<string, string | undefined>) =>
              `http://localhost:8081/api/conversations/export?status=${filters['status'] ?? ''}`,
            createConversation: createConversationSpy,
          },
        },
      ],
    }).compileComponents();
  });

  it('should display conversations returned by the service, including tags', () => {
    const fixture = TestBed.createComponent(Conversations);

    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;

    expect(page.textContent).toContain('John Smith');
    expect(page.textContent).toContain('Unable to reset password');
    expect(page.textContent).toContain('Open');
    expect(page.textContent).toContain('High');
    expect(page.textContent).toContain('VIP');
    expect(page.textContent).toContain('Login');
    expect(page.textContent).toContain('Alex Chen');

    const exportLink = page.querySelector('.export-link') as HTMLAnchorElement;
    expect(exportLink.getAttribute('href')).toContain('/api/conversations/export');
  });

  it('should show conversation details, including notes and follow-up date, when selected', async () => {
    const fixture = TestBed.createComponent(Conversations);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    const card = page.querySelector('.conversation-card') as HTMLButtonElement;
    card.click();
    fixture.detectChanges();
    await fixture.whenStable();

    expect(page.textContent).toContain('john.smith@example.com');

    const notesField = page.querySelector('#edit-notes') as HTMLTextAreaElement;
    expect(notesField.value).toBe('Sent a temporary reset link.');

    const followUpInput = page.querySelector('#edit-follow-up') as HTMLInputElement;
    expect(followUpInput.value).toBe('2026-09-29');

    const assigneeSelect = page.querySelector('#edit-assignee') as HTMLSelectElement;
    expect(assigneeSelect.value).toBe('Alex Chen');

    expect(page.textContent).toContain('AI Insights');
    expect(page.textContent).toContain(mockInsights.summary);
    expect(page.textContent).toContain('Matches current priority');
  });

  it('creates a conversation and selects it in the refreshed list', () => {
    const fixture = TestBed.createComponent(Conversations);
    fixture.detectChanges();

    const component = fixture.componentInstance;
    component.openCreateForm();
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    expect(page.querySelector('.modal-backdrop')).toBeTruthy();

    component.createConversation({
      customerName: 'Jane Doe',
      customerEmail: 'jane@example.com',
      subject: 'Cannot log in',
      priority: 'MEDIUM',
    });
    fixture.detectChanges();

    expect(createConversationSpy).toHaveBeenCalledWith({
      customerName: 'Jane Doe',
      customerEmail: 'jane@example.com',
      subject: 'Cannot log in',
      priority: 'MEDIUM',
    });
    expect(component.showCreateForm()).toBe(false);
    expect(component.selectedConversation()?.id).toBe('2');
    expect(page.querySelector('.modal-backdrop')).toBeFalsy();
  });

  it('keeps the create form open and shows an error when creation fails', () => {
    createConversationSpy.mockReturnValue(throwError(() => new Error('boom')));

    const fixture = TestBed.createComponent(Conversations);
    fixture.detectChanges();

    const component = fixture.componentInstance;
    component.openCreateForm();
    component.createConversation({
      customerName: '',
      customerEmail: '',
      subject: '',
    });
    fixture.detectChanges();

    expect(component.showCreateForm()).toBe(true);
    expect(component.createError()).toContain('Failed to create conversation');
  });
});
