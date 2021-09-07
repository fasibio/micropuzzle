
import type { ComponentType } from 'react'
import React from 'react'
import ReactDOM from 'react-dom'
import { StyleSheetManager } from 'styled-components'

type BaseProps = { [key: string]: string | undefined }
type Args<P extends BaseProps> = {
    name: string
    attributes: string[]
    component: ComponentType<P>
}

/**
 * Register a custom element that wraps a React component.
 *
 * @param name       - the name of the custom element
 * @param attributes - the names of the custom element's attributes
 * @param component  - the React component
 */
export function registerCustomElement<P extends BaseProps>({
    name,
    attributes,
    component: Component,
}: Args<P>) {
    const webComponentClass = class extends HTMLElement {
        private readonly styleHost: HTMLElement
        private readonly mountPoint: HTMLElement

        constructor() {
            super()
            this.styleHost = document.createElement('div')
            this.mountPoint = document.createElement('div')
            this.attachShadow({ mode: 'open' })
        }

        connectedCallback() {
            if (this.isConnected) {
                const attrs = attributes.reduce(
                    (acc, key) =>
                        Object.assign(acc, {
                            [key]: this.getAttribute(key) ?? undefined,
                        }),
                    {} as P
                )

                this.shadowRoot?.appendChild(this.styleHost)
                this.shadowRoot?.appendChild(this.mountPoint)
                ReactDOM.render(
                    <StyleSheetManager target={this.styleHost}>
                        <Component {...attrs} />
                    </StyleSheetManager>,
                    this.mountPoint
                )
            }
        }

        disconnectedCallback() {
            if (!this.isConnected) {
                this.shadowRoot?.removeChild(this.mountPoint)
                this.shadowRoot?.removeChild(this.styleHost)
            }
        }
    }

    customElements.define(name, webComponentClass)
}
