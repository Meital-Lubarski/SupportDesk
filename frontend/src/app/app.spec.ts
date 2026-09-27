import { TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { of } from 'rxjs';
import { App } from './app';
import { DashboardSummary } from './conversation';
import { ConversationService } from './conversation.service';
import { Home } from './home/home';

describe('App', () => {
  const mockSummary: DashboardSummary = {
    total: 0,
    byStatus: {},
    byPriority: {},
    byTag: {},
    byAgent: {},
    overdueFollowUps: 0,
    upcomingFollowUps: 0,
    suggestedPriorityBumps: 0,
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [App],
      providers: [
        {
          provide: ConversationService,
          useValue: {
            getConversations: () => of([]),
            getDashboardSummary: () => of(mockSummary),
            getAgents: () => of([]),
            getExportUrl: () => 'http://localhost:8081/api/conversations/export',
          },
        },
      ],
    }).compileComponents();
  });

  it('shows the home tab by default', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    expect(page.querySelector('app-home')).toBeTruthy();
    expect(page.querySelector('app-conversations')).toBeFalsy();
    expect(page.querySelector('app-dashboard')).toBeFalsy();
  });

  it('switches to the dashboard tab when clicked', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    const buttons = Array.from(page.querySelectorAll('button.tab'));
    const dashboardTab = buttons.find((b) => b.textContent?.trim() === 'Dashboard') as HTMLButtonElement;

    dashboardTab.click();
    fixture.detectChanges();

    expect(page.querySelector('app-dashboard')).toBeTruthy();
    expect(page.querySelector('app-home')).toBeFalsy();
  });

  it('navigates to conversations when the home page CTA is used', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    const page = fixture.nativeElement as HTMLElement;
    const home = fixture.debugElement.query(By.directive(Home)).componentInstance as Home;
    home.goTo('conversations');
    fixture.detectChanges();

    expect(page.querySelector('app-conversations')).toBeTruthy();
    expect(page.querySelector('app-home')).toBeFalsy();
  });
});
