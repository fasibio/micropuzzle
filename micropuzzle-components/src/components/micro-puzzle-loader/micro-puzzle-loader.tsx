import { Component, Host, h, Element, Event, EventEmitter } from '@stencil/core';

import {NewContentEventDetails} from '../../utils/utils'


@Component({
  tag: 'micro-puzzle-loader',
  shadow: false,
})
export class MicroPuzzleLoader {
  @Element() el: HTMLElement;
  @Event({
    eventName: 'new-content',
    bubbles: true,
    composed: true,
    cancelable: true
  }) newContentEvent: EventEmitter<NewContentEventDetails>

  constructor(){
    setInterval(() => {
      console.log('send event')
      this.newContentEvent.emit({
        content: "<h2>lalal</h2>",
        name: "footer"
      })
    }, 3000)
  }

  render() {
    return (
      <Host />
    );
  }

}
