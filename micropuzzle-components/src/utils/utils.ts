export const NEW_CONTENT_EVENT = 'new-content';

export type NewContentEventDetails = {
  content: string;
  name: string;
};

export type LoadContentPayload = {
  content: string;
  loading: string;
};
