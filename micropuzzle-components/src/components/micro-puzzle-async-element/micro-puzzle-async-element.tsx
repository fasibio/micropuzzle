import { Component, Host, h, Element } from '@stencil/core';

@Component({
  tag: 'micro-puzzle-async-element',
  shadow: true,
})
export class MicroPuzzleAsyncElement {



  render() {
    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }

}
