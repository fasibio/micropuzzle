import { Component, Host, h, Element, Event, EventEmitter } from '@stencil/core';

import {NewContentEventDetails} from '../../utils/utils'


@Component({
  tag: 'micro-puzzle-loader',
  shadow: false,
})
export class MicroPuzzleLoader {
  @Element() el: HTMLElement;
  @Event({
    eventName: "test1234",
    bubbles: true,
    composed: true,
    cancelable: true
  }) testEvent: EventEmitter<NewContentEventDetails>

  constructor(){
    setInterval(() => {
      console.log('send event')
      this.testEvent.emit({
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
