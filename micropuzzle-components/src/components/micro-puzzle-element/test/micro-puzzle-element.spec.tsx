import { newSpecPage } from '@stencil/core/testing';
import { MicroPuzzleElement } from '../micro-puzzle-element';

describe('micro-puzzle-element', () => {
  it('renders', async () => {
    const page = await newSpecPage({
      components: [MicroPuzzleElement],
      html: `<micro-puzzle-element></micro-puzzle-element>`,
    });
    expect(page.root).toEqualHtml(`
      <micro-puzzle-element>
        <mock:shadow-root>
          <slot></slot>
        </mock:shadow-root>
      </micro-puzzle-element>
    `);
  });
});
