export type FollowUpState = 'none' | 'overdue' | 'soon' | 'scheduled';

//Classifies a follow-up date relative to today for badge styling.
export function followUpState(
  followUpDate: string,
  status: string,
): FollowUpState {
  if (!followUpDate || status === 'RESOLVED') {
    return 'none';
  }

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const target = new Date(followUpDate);
  const diffDays = (target.getTime() - today.getTime()) / 86_400_000;

  if (diffDays < 0) {
    return 'overdue';
  }

  if (diffDays <= 7) {
    return 'soon';
  }

  return 'scheduled';
}
