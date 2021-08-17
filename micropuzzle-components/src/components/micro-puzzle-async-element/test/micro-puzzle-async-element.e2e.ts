import { newE2EPage } from '@stencil/core/testing';

describe('micro-puzzle-async-element', () => {
  it('renders', async () => {
    const page = await newE2EPage();
    await page.setContent('<micro-puzzle-async-element></micro-puzzle-async-element>');

    const element = await page.find('micro-puzzle-async-element');
    expect(element).toHaveClass('hydrated');
  });
});
