import { newE2EPage } from '@stencil/core/testing';

describe('micro-puzzle-loader', () => {
  it('renders', async () => {
    const page = await newE2EPage();
    await page.setContent('<micro-puzzle-loader></micro-puzzle-loader>');

    const element = await page.find('micro-puzzle-loader');
    expect(element).toHaveClass('hydrated');
  });
});
