import { Component, Host, h, Element, Event, EventEmitter, Prop, Listen } from '@stencil/core';

import {NewContentEventDetails, LoadContentPayload} from '../../utils/utils'
import io from "socket.io-client";


@Component({
  tag: 'micro-puzzle-loader',
  shadow: false,
})
export class MicroPuzzleLoader {
  @Element() el: HTMLMicroPuzzleLoaderElement;
  /**
   * This is a internal micropuzzle event. If new Content is there to show it 
   */
  @Event({
    eventName: 'new-content',
    bubbles: true,
    composed: true,
    cancelable: true
  }) newContentEvent: EventEmitter<NewContentEventDetails>

  /**
   * The UUID given by the Micropuzzle SSI side
   */
  @Prop() streamregistername: string

  /**
   * The URL where to connect the streamingconnection
   */
  @Prop() streamingurl: string

  private socket: SocketIOClient.Socket;
  constructor(){
    this.socket = io({
      query: `streamId=${this.streamregistername}`
    })
    this.socket.on("NEW_CONTENT", (data: {key: string, value: string}) => {
      this.newContentEvent.emit({
        content: data.value,
        name:data.key
      })
    })
   
  }

  @Listen('load-content', {target: "window"})
  loadNewContent(event: CustomEvent<LoadContentPayload>){
    this.socket.emit("LOAD_CONTENT", event.detail)
  }

  render() {
    console.log('da', this.streamregistername)
    return (
      <Host >
      </Host>
    );
  }

}
