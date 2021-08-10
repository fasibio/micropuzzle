import { Component, Host, h, Element } from '@stencil/core';

@Component({
  tag: 'micro-puzzle-loader',
  shadow: true,
})
export class MicroPuzzleLoader {
  @Element() el: HTMLElement;

  constructor(){
    setInterval(() => {
      console.log('send event')
      this.el.dispatchEvent(new CustomEvent("TEST", 
      {
        bubbles: true, 
        composed: true, 
        detail: {content: "<span>test</span>"}}))
    }, 1000)
  }

  render() {
    return (
      <Host />
    );
  }

}
