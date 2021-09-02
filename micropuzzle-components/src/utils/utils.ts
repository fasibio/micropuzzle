export const NEW_CONTENT_EVENT = 'new-content';

export type NewContentEventDetails = {
  content: string;
  name: string;
};

export type LoadContentPayload = {
  content: string;
  loading: string;
};

export const sleep = (timeMs: number): Promise<unknown> => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve('');
    }, timeMs);
  });
};
