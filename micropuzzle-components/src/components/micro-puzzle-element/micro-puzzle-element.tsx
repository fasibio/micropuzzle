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
  element: ShadowRoot


  @Listen('new-content', {target: "window"})
  eventUpdated(event: CustomEvent<NewContentEventDetails>){
    if (event.detail.name === this.name){
      this.element.innerHTML =event.detail.content 
      event.preventDefault()
    }
  }

  componentDidRender(){
    const result = this.el.getElementsByTagName("template")[0]
    const templateContent = result.content
    this.element = this.el.attachShadow({mode: 'open'})
    this.element.appendChild(templateContent.cloneNode(true))
  }

  render() {

    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }

}
