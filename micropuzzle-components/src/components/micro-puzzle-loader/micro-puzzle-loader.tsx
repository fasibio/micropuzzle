import { Component, Host, h, Element, Event, EventEmitter, Prop } from '@stencil/core';

import {NewContentEventDetails} from '../../utils/utils'
import io from "socket.io-client";


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

  @Prop() streamregistername: string
  @Prop() streamingurl: string

  private socket: SocketIOClient.Socket;
  constructor(){
    this.socket =  io({
      query: `streamId=${this.streamregistername}`
    })
    this.socket.on("NEW_CONTENT", (data: {key: string, value: string}) => {
      this.newContentEvent.emit({
        content: data.value,
        name:data.key
      })
    })
   
  }

  render() {
    console.log('da', this.streamregistername)
    return (
      <Host >
      </Host>
    );
  }

}
