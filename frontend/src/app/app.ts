import { Component, signal } from '@angular/core';
import { Conversations } from './conversations/conversations';
import { Dashboard } from './dashboard/dashboard';
import { Home } from './home/home';
import { AppView } from './shared/app-view';

//Root shell: holds the active tab.
@Component({
  selector: 'app-root',
  imports: [Home, Conversations, Dashboard],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  activeView = signal<AppView>('home');

  setView(view: AppView): void {
    this.activeView.set(view);
  }
}
