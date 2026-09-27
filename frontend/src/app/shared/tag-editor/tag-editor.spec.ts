import { TestBed } from '@angular/core/testing';
import { TagEditor } from './tag-editor';

describe('TagEditor', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TagEditor],
    }).compileComponents();
  });

  it('emits a new tag list when a tag is added', () => {
    const fixture = TestBed.createComponent(TagEditor);
    fixture.componentRef.setInput('tags', ['Billing']);
    fixture.detectChanges();

    let emitted: string[] | undefined;
    fixture.componentInstance.tagsChange.subscribe((tags) => (emitted = tags));

    fixture.componentInstance.newTag.set('VIP');
    fixture.componentInstance.addTag();

    expect(emitted).toEqual(['Billing', 'VIP']);
  });

  it('does not add a duplicate tag (case-insensitive)', () => {
    const fixture = TestBed.createComponent(TagEditor);
    fixture.componentRef.setInput('tags', ['Billing']);
    fixture.detectChanges();

    let emitted: string[] | undefined;
    fixture.componentInstance.tagsChange.subscribe((tags) => (emitted = tags));

    fixture.componentInstance.newTag.set('billing');
    fixture.componentInstance.addTag();

    expect(emitted).toBeUndefined();
  });

  it('emits a new tag list when a tag is removed', () => {
    const fixture = TestBed.createComponent(TagEditor);
    fixture.componentRef.setInput('tags', ['Billing', 'VIP']);
    fixture.detectChanges();

    let emitted: string[] | undefined;
    fixture.componentInstance.tagsChange.subscribe((tags) => (emitted = tags));

    fixture.componentInstance.removeTag('Billing');

    expect(emitted).toEqual(['VIP']);
  });
});
