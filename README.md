# Support Desk

A full-stack customer support conversation workspace: search, tag, assign, and follow up on conversations, back it with a live dashboard, and get a rule-based "AI" assist on what to prioritize next.

- **Frontend:** Angular (standalone components, signals) + TypeScript
- **Backend:** Go, `net/http`, no framework
- **Storage:** in-memory (see [Limitations](#limitations))

## Features

**Home** — a landing screen with a live dashboard snapshot (total conversations, overdue follow-ups, AI-flagged priority bumps) and quick entry points into Conversations and Dashboard.

**Conversations**
- Create a conversation from a modal form (customer name/email/subject required; priority, assignee, tags optional)
- Search by customer name, email, or subject
- Filter by status, priority, tag, and assignee (including an "Unassigned" filter)
- Click a conversation to expand an inline edit form: status, priority, follow-up date, assignee, tags, notes
- Assign conversations to a support agent from a fixed roster
- Export the currently filtered list as a CSV file
- Loading, error, and empty states throughout

**AI Insights** (per conversation) — a synthesized one-line summary and a suggested priority with a one-click "Apply" button when it differs from the current one. **This is not a call to a real language model** — there's no API key wired up. It's a deterministic keyword/heuristic scoring engine (urgency keywords, VIP tag, overdue follow-ups, days-open) plus a templated sentence builder, fully unit tested, and the UI says so ("rule-based, not a live model"). Swapping in a real LLM later is a self-contained change to `suggestPriority`/`summarize` in `backend/insights.go`.

**Activity timeline** (per conversation) — an automatic audit trail. Every edit is diffed against the previous state and recorded as a human-readable entry ("Status changed from open to resolved", "Reassigned from Alex Chen to Priya Patel") — no manual logging, and no-op saves don't spam the log.

**Dashboard**
- Stat tiles: total conversations, overdue follow-ups, follow-ups due in 7 days, AI-flagged priority bumps
- Status distribution: a stacked bar chart (Open → In progress → Resolved) with a legend and hover tooltips
- Breakdowns by priority, by tag, and by agent workload
- "Last updated" timestamp with a manual refresh that doesn't blank the page while it reloads

## API reference

All endpoints are served from `http://localhost:8081/api`.

| Method | Path | Description |
|---|---|---|
| `GET` | `/conversations` | List conversations. Query params: `status`, `priority`, `tag`, `assignee` (or `unassigned`), `search` |
| `POST` | `/conversations` | Create a conversation (`customerName`, `customerEmail`, `subject` required) |
| `GET` | `/conversations/export` | Same filters as above, returns a CSV file |
| `GET` | `/conversations/{id}` | Get one conversation |
| `PATCH` | `/conversations/{id}` | Update status, priority, tags, followUpDate, notes, and/or assignedTo |
| `GET` | `/conversations/{id}/insights` | AI summary + suggested priority for one conversation |
| `GET` | `/conversations/{id}/activity` | Change history for one conversation |
| `GET` | `/agents` | The fixed support-agent roster |
| `GET` | `/dashboard` | Aggregated counts for the dashboard |

## Project structure

```
backend/
  main.go            routes
  conversation.go    data model + seed data
  handlers.go         list/create/update conversations, filtering
  agents.go            agent roster + validation
  insights.go          AI-ish heuristic engine
  activity.go           change-diffing + audit log
  dashboard.go           aggregation for the dashboard
  export.go               CSV export
  middleware.go            CORS
  *_test.go                one test file per capability above

frontend/src/app/
  home/                       landing page
  conversations/               list, filters, orchestration
  conversation-detail/          the edit form (status/priority/tags/notes/assignee)
  new-conversation/              the "create conversation" modal
  dashboard/                      stat tiles, chart, breakdowns
  shared/
    tag-editor/                    reusable tag add/remove UI
    badge.util.ts                  status/priority → color + label (one source of truth)
    follow-up.util.ts               date → overdue/soon/scheduled
    app-view.ts                      the app's tab-navigation type
```

## Getting started

### Backend

```bash
cd backend
go run .
```

Runs at `http://localhost:8081`.

### Frontend

```bash
cd frontend
npm install
npm start
```

Runs at `http://localhost:4200`. On Windows PowerShell, if `npm.ps1` execution is blocked, use `npm.cmd install` / `npm.cmd start`.

## Testing

```bash
cd backend && go test -v ./...      # 25 tests
cd frontend && npm test -- --watch=false   # 20 tests
```

Business logic (priority scoring, activity diffing, dashboard aggregation, CSV filtering) is written as pure functions that take an injected `time.Now()`, so it's tested directly without spinning up HTTP servers or mocking the clock.

## Limitations

- **No concurrency safety.** `conversations` and `activityLog` are plain in-memory maps/slices with no mutex. Fine for one local process with one user; not safe for concurrent traffic. A real deployment would need a real database (or at least a `sync.Mutex`).
- **Data doesn't persist.** Restarting the backend resets everything to the seed data.
- **AI Insights are heuristic, not generative** — see the Features section above.
- **Agent roster is hardcoded** (`backend/agents.go`) rather than backed by a user directory.

## Project history

This started as a timed take-home assignment (a minimal Angular + Go conversation viewer) and was substantially redesigned and extended afterward: visual redesign, componentization, and the Home/AI Insights/Activity/Dashboard chart/agent-assignment/CSV-export/create-conversation features were all added in a later pass.
