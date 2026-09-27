//Maps a status or priority value to a CSS badge class.
export function statusBadgeClass(status: string): string {
  switch (status) {
    case 'OPEN':
      return 'badge badge-status-open';
    case 'IN_PROGRESS':
      return 'badge badge-status-in-progress';
    case 'RESOLVED':
      return 'badge badge-status-resolved';
    default:
      return 'badge';
  }
}

export function priorityBadgeClass(priority: string): string {
  switch (priority) {
    case 'LOW':
      return 'badge badge-priority-low';
    case 'MEDIUM':
      return 'badge badge-priority-medium';
    case 'HIGH':
      return 'badge badge-priority-high';
    default:
      return 'badge';
  }
}

export function statusLabel(status: string): string {
  switch (status) {
    case 'OPEN':
      return 'Open';
    case 'IN_PROGRESS':
      return 'In progress';
    case 'RESOLVED':
      return 'Resolved';
    default:
      return status;
  }
}

export function priorityLabel(priority: string): string {
  return priority.charAt(0) + priority.slice(1).toLowerCase();
}
