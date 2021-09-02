import { Component, Host, h, Element, Event, EventEmitter, Prop, Listen } from '@stencil/core';
import { NewContentEventDetails, LoadContentPayload } from '../../utils/utils';

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
    cancelable: true,
  })
  newContentEvent: EventEmitter<NewContentEventDetails>;

  /**
   * The UUID given by the Micropuzzle SSI side
   */
  @Prop() streamregistername: string;

  /**
   * The URL where to connect the streamingconnection
   */
  @Prop() streamingurl: string;

  private socket: WebSocket;
  constructor() {
    this.socket = new WebSocket(`${this.streamingurl}?streamid=${this.streamregistername}`);
    this.socket.onmessage = event => {
      const data = JSON.parse(event.data) as {
        type: string;
        data: unknown;
      };
      if (data.type === 'NEW_CONTENT') {
        const newContentData = data.data as {
          key: string;
          value: string;
        };
        this.newContentEvent.emit({
          content: newContentData.value,
          name: newContentData.key,
        });
      }
    };
    // this.socket.on("NEW_CONTENT", (data: {key: string, value: string}) => {
    //   console.log(data.value)
    //   this.newContentEvent.emit({
    //     content: data.value,
    //     name:data.key
    //   })
    // })
  }

  @Listen('load-content', { target: 'window' })
  loadNewContent(event: CustomEvent<LoadContentPayload>) {
    this.socket.send(
      JSON.stringify({
        type: 'LOAD_CONTENT',
        data: event.detail,
      }),
    );
  }

  render() {
    return <Host></Host>;
  }
}
