export const NEW_CONTENT_EVENT = 'new-content';

export type NewContentEventDetails = {
  content: string;
  name: string;
};

export type LoadContentPayload = {
  content: string;
  loading: string;
};

export interface NewFragmentPayload {
  key: string;
  value: string;
  isFallback: boolean;
}

export const sleep = (timeMs: number): Promise<unknown> => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve('');
    }, timeMs);
  });
};
