import {
  Component,
  EventEmitter,
  HostListener,
  Input,
  Output,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CreateConversationRequest } from '../conversation';
import { TagEditor } from '../shared/tag-editor/tag-editor';

//Modal form for creating a new conversation. A "dumb" form like
//ConversationDetail: it emits the request and lets the parent own the
//service call, loading state, and error message.
@Component({
  selector: 'app-new-conversation',
  imports: [FormsModule, TagEditor],
  templateUrl: './new-conversation.html',
  styleUrl: './new-conversation.css',
})
export class NewConversation {
  @Input() agents: string[] = [];
  @Input() creating = false;
  @Input() errorMessage = '';

  @Output() create = new EventEmitter<CreateConversationRequest>();
  @Output() cancel = new EventEmitter<void>();

  customerName = '';
  customerEmail = '';
  subject = '';
  priority = 'MEDIUM';
  assignedTo = '';
  tags: string[] = [];

  @HostListener('document:keydown.escape')
  onEscape(): void {
    this.cancel.emit();
  }

  submit(): void {
    this.create.emit({
      customerName: this.customerName,
      customerEmail: this.customerEmail,
      subject: this.subject,
      priority: this.priority,
      assignedTo: this.assignedTo,
      tags: this.tags,
    });
  }
}
