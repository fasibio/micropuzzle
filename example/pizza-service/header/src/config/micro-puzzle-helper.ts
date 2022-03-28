/**
 * Mircopuzzle AUTO-GENERATED CODE: PLEASE DO NOT MODIFY MANUALLY
 */
	
export enum MicropuzzleFrontends {
	CART_CONTENT="cart.content",
	DETAIL_CONTENT="detail.content",
	GLOBAL_FALLBACK="global.fallback",
	GLOBAL_FOOTER="global.footer",
	GLOBAL_HEADER="global.header",
	STARTPAGE_CONTENT="startpage.content",
}

export const loadMicroFrontend = (area: string, microfrontend: MicropuzzleFrontends ) =>{
  const event = new CustomEvent('load-content',{
    bubbles: true,
    composed: true,
    cancelable: true,
    detail: {
    loading: microfrontend,
    content: area
    }
  })
  document.dispatchEvent(event);
}