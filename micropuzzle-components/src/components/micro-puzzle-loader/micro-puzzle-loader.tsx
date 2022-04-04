import { Component, Host, h, Element, Event, EventEmitter, Prop, Listen } from '@stencil/core';
import { NewContentEventDetails, NewFragmentPayload, LoadContentPayload, PageDeclarations } from '../../utils/utils';

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

  /**
   * The count of Fallbacks by initial loading
   */
  @Prop() fallbacks: number;

  /**
   * Declarations of the pages
   */
  @Prop() pagesstr: string

  private waitingLoadings = 0;
  private asyncLoadingCount = 0;
  private socket: WebSocket;
  constructor() {
    const page :PageDeclarations= JSON.parse(this.pagesstr) ;
    window.onpopstate = ( event: PopStateEvent & {target: {location: {pathname: string}}}) => {
      const currentPage = Object.values(page).find(onePage => {
        return onePage.url === event.target.location.pathname;
      });
      Object.keys(currentPage.fragments).forEach(fragment => {
        this.loadNewContent({detail: {
          content: fragment,
          loading: currentPage.fragments[fragment]
        }} as any);
      });
    }
    

    this.waitingLoadings = this.fallbacks;
    if (this.fallbacks > 0) {
      this.startSocketConnection();
    }
  }

  private getSocketUrl(): string {
    const loc = window.location;
    let new_uri = 'ws:';
    if (loc.protocol === 'https:') {
      new_uri = 'wss:';
    }
    new_uri += '//' + loc.host;
    new_uri += '/' + this.streamingurl;
    return new_uri;
  }

  private startSocketConnection() {
    if (this.socket === undefined || this.socket.readyState !== WebSocket.OPEN) {
      this.socket = new WebSocket(`${this.getSocketUrl()}?streamid=${this.streamregistername}`);
      this.socket.onmessage = event => {
        const data = JSON.parse(event.data) as {
          type: string;
          data: unknown;
        };
        if (data.type === 'NEW_CONTENT') {
          this.sendNewContentEvent(data.data as NewFragmentPayload);
          this.asyncLoadingCount = this.asyncLoadingCount + 1;
          if (this.waitingLoadings === this.asyncLoadingCount) {
            this.waitingLoadings = 0;
            this.asyncLoadingCount = 0;
            this.socket.close();
          }
        }
      };
    }
  }

  private sendNewContentEvent(data: NewFragmentPayload) {
    this.newContentEvent.emit({
      content: data.value,
      name: data.key,
    });
  }

  @Listen('load-content', { target: 'window' })
  async loadNewContent(event: CustomEvent<LoadContentPayload>) {
    const fetchResult = await fetch(`/micro-puzzle?fragment=${event.detail.loading}&frontend=${event.detail.content}`, {headers: {
      streamid: this.streamregistername
    }});
    const data: NewFragmentPayload = await fetchResult.json();
    this.sendNewContentEvent(data);
    if (data.isFallback) {
      this.startSocketConnection();
      this.waitingLoadings = this.waitingLoadings + 1;
    }
  }

  render() {
    return <Host></Host>;
  }
}

