/* eslint-disable */
/* tslint:disable */
/**
 * This is an autogenerated file created by the Stencil compiler.
 * It contains typing information for all components that exist in this project.
 */
import { HTMLStencilElement, JSXBase } from "@stencil/core/internal";
export namespace Components {
    interface MicroPuzzleElement {
        "name": string;
    }
}
declare global {
    interface HTMLMicroPuzzleElementElement extends Components.MicroPuzzleElement, HTMLStencilElement {
    }
    var HTMLMicroPuzzleElementElement: {
        prototype: HTMLMicroPuzzleElementElement;
        new (): HTMLMicroPuzzleElementElement;
    };
    interface HTMLElementTagNameMap {
        "micro-puzzle-element": HTMLMicroPuzzleElementElement;
    }
}
declare namespace LocalJSX {
    interface MicroPuzzleElement {
        "name"?: string;
    }
    interface IntrinsicElements {
        "micro-puzzle-element": MicroPuzzleElement;
    }
}
export { LocalJSX as JSX };
declare module "@stencil/core" {
    export namespace JSX {
        interface IntrinsicElements {
            "micro-puzzle-element": LocalJSX.MicroPuzzleElement & JSXBase.HTMLAttributes<HTMLMicroPuzzleElementElement>;
        }
    }
}