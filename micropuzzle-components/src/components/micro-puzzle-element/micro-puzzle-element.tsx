import { Component, Host, h, Prop, Element,  } from '@stencil/core';

@Component({
  tag: 'micro-puzzle-element',
  styleUrl: 'micro-puzzle-element.css',
  shadow: false,
})
export class MicroPuzzleElement {
  @Element() el: HTMLElement;
  @Prop() name: string

  componentDidRender(){
    const result = this.el.getElementsByTagName("template")[0]
    const templateContent = result.content
    this.el.attachShadow({mode: 'open'}).appendChild(templateContent.cloneNode(true))
  }

  render() {

    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }

}
