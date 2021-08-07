import { Component, Host, h, Prop } from '@stencil/core';

@Component({
  tag: 'micro-puzzle-element',
  styleUrl: 'micro-puzzle-element.css',
  shadow: true,
})
export class MicroPuzzleElement {

  @Prop() name: string

  render() {
    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }

}
