import { Component, Host, h, Prop, Element, Listen, State } from '@stencil/core';
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
      const result = this.el.getElementsByTagName('template')[0];
      console.log(result.content);
      result.innerHTML = '';
      const p = new DOMParser().parseFromString(event.detail.content, 'text/html');
      
      const scripts = p.getElementsByTagName('script')
      for (let i = 0; i < scripts.length; i++) {
        const script = document.createElement('script');
        scripts[i].getAttributeNames().forEach(attr => {
          script.setAttribute(attr, scripts[i].getAttribute(attr));
        })
        script.textContent = scripts[i].textContent;
        script.async = false;
        result.content.appendChild(script);
      }

      while (scripts.length > 0) {
        const script = document.createElement('script');
        scripts[0].getAttributeNames().forEach(attr => {
          script.setAttribute(attr, scripts[0].getAttribute(attr));
        })
        script.textContent = scripts[0].textContent;
        script.async = false;
        result.content.appendChild(script);
        scripts[0].remove();
      }
      const headElements = p.querySelectorAll('head > *')
      for (let i = 0; i < headElements.length; i++) {
        result.content.appendChild(headElements[i].cloneNode(true));
      }
      const bodyElements = p.querySelectorAll('body > *')
      for (let i = 0; i < bodyElements.length; i++) {
        result.content.appendChild(bodyElements[i].cloneNode(true));
      }
      const templateContent = result.content;
      this.element.innerHTML = ''
      this.element.appendChild(templateContent.cloneNode(true));
    }
  }

  componentDidRender() {
    const result = this.el.getElementsByTagName('template')[0];
    const templateContent = result.content;
    if (!this.el.shadowRoot){
      this.element = this.el.attachShadow({ mode: 'open' });
    }
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
