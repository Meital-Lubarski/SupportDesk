import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';

//Add/remove chips for a tag list. Shared by the conversation detail form and
//the new-conversation form so tag-editing logic lives in exactly one place.
@Component({
  selector: 'app-tag-editor',
  imports: [FormsModule],
  templateUrl: './tag-editor.html',
  styleUrl: './tag-editor.css',
})
export class TagEditor {
  @Input() tags: string[] = [];
  @Input() inputId = 'tag-editor-input';
  @Output() tagsChange = new EventEmitter<string[]>();

  newTag = signal('');

  addTag(): void {
    const tag = this.newTag().trim();

    if (!tag || this.tags.some((t) => t.toLowerCase() === tag.toLowerCase())) {
      this.newTag.set('');
      return;
    }

    this.tagsChange.emit([...this.tags, tag]);
    this.newTag.set('');
  }

  removeTag(tag: string): void {
    this.tagsChange.emit(this.tags.filter((t) => t !== tag));
  }
}
