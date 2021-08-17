import { newSpecPage } from '@stencil/core/testing';
import { MicroPuzzleAsyncElement } from '../micro-puzzle-async-element';

describe('micro-puzzle-async-element', () => {
  it('renders', async () => {
    const page = await newSpecPage({
      components: [MicroPuzzleAsyncElement],
      html: `<micro-puzzle-async-element></micro-puzzle-async-element>`,
    });
    expect(page.root).toEqualHtml(`
      <micro-puzzle-async-element>
        <mock:shadow-root>
          <slot></slot>
        </mock:shadow-root>
      </micro-puzzle-async-element>
    `);
  });
});
