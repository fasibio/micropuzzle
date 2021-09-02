import { Component, Host, h, Prop, Element, Listen } from '@stencil/core';
import { NewContentEventDetails, sleep } from '../../utils/utils';

@Component({
  tag: 'micro-puzzle-element',
  shadow: false,
})
export class MicroPuzzleElement {
  @Element() el: HTMLMicroPuzzleElementElement;

  /**
   * The Logic unique name for this elementarea
   */
  @Prop() name: string;

  private element: ShadowRoot;

  @Listen('new-content', { target: 'window' })
  async eventUpdated(event: CustomEvent<NewContentEventDetails>) {
    if (event.detail.name === this.name) {
      while (this.element === undefined) {
        await sleep(1);
      }
      this.element.innerHTML = event.detail.content;
      event.preventDefault();
    }
  }

  componentDidRender() {
    const result = this.el.getElementsByTagName('template')[0];
    const templateContent = result.content;
    this.element = this.el.attachShadow({ mode: 'open' });
    this.element.appendChild(templateContent.cloneNode(true));
  }

  render() {
    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }
}
