export function format(first: string, middle: string, last: string): string {
  return (first || '') + (middle ? ` ${middle}` : '') + (last ? ` ${last}` : '');
}

export const NEW_CONTENT_EVENT = 'new-content';

export type NewContentEventDetails = {
  content: string;
  name: string;
};
