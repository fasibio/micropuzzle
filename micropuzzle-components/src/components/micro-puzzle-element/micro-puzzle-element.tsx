import { Component, Host, h, Prop, Element, Listen } from '@stencil/core';
import {NewContentEventDetails} from '../../utils/utils'

@Component({
  tag: 'micro-puzzle-element',
  styleUrl: 'micro-puzzle-element.css',
  shadow: false,
})
export class MicroPuzzleElement {
  @Element() el: HTMLElement;
  @Prop() name: string

  @Listen('new-content', {target: "window"})
  eventUpdated(event: CustomEvent<NewContentEventDetails>){
    if (event.detail.name === this.name){
      console.log('hier', event.detail.content)
      event.preventDefault()
    } else {
      console.log('event nicht für mich')
    }

  }

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
