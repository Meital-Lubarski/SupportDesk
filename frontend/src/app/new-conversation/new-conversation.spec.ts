import { TestBed } from '@angular/core/testing';
import { CreateConversationRequest } from '../conversation';
import { NewConversation } from './new-conversation';

describe('NewConversation', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NewConversation],
    }).compileComponents();
  });

  it('emits the form values on submit', () => {
    const fixture = TestBed.createComponent(NewConversation);
    fixture.detectChanges();

    const component = fixture.componentInstance;
    component.customerName = 'Jane Doe';
    component.customerEmail = 'jane@example.com';
    component.subject = 'Cannot log in';
    component.priority = 'HIGH';
    component.assignedTo = 'Alex Chen';
    component.tags = ['VIP'];

    let emitted: CreateConversationRequest | undefined;
    component.create.subscribe((request) => (emitted = request));

    component.submit();

    expect(emitted).toEqual({
      customerName: 'Jane Doe',
      customerEmail: 'jane@example.com',
      subject: 'Cannot log in',
      priority: 'HIGH',
      assignedTo: 'Alex Chen',
      tags: ['VIP'],
    });
  });

  it('emits cancel when the backdrop is clicked', () => {
    const fixture = TestBed.createComponent(NewConversation);
    fixture.detectChanges();

    let cancelled = false;
    fixture.componentInstance.cancel.subscribe(() => (cancelled = true));

    const backdrop = fixture.nativeElement.querySelector('.modal-backdrop') as HTMLElement;
    backdrop.click();

    expect(cancelled).toBe(true);
  });

  it('emits cancel on Escape', () => {
    const fixture = TestBed.createComponent(NewConversation);
    fixture.detectChanges();

    let cancelled = false;
    fixture.componentInstance.cancel.subscribe(() => (cancelled = true));

    fixture.componentInstance.onEscape();

    expect(cancelled).toBe(true);
  });
});
