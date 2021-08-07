import { newE2EPage } from '@stencil/core/testing';

describe('micro-puzzle-element', () => {
  it('renders', async () => {
    const page = await newE2EPage();
    await page.setContent('<micro-puzzle-element></micro-puzzle-element>');

    const element = await page.find('micro-puzzle-element');
    expect(element).toHaveClass('hydrated');
  });
});
