import { newSpecPage } from '@stencil/core/testing';
import { MicroPuzzleLoader } from '../micro-puzzle-loader';

describe('micro-puzzle-loader', () => {
  it('renders', async () => {
    const page = await newSpecPage({
      components: [MicroPuzzleLoader],
      html: `<micro-puzzle-loader></micro-puzzle-loader>`,
    });
    expect(page.root).toEqualHtml(`
      <micro-puzzle-loader>
        <mock:shadow-root>
          <slot></slot>
        </mock:shadow-root>
      </micro-puzzle-loader>
    `);
  });
});
