/**
 * Mircopuzzle AUTO-GENERATED CODE: PLEASE DO NOT MODIFY MANUALLY
 */
	
export enum MicropuzzleFrontends {
	ABOUT_CONTENT="about.content",
	GLOBAL_FALLBACK="global.fallback",
	GLOBAL_FOOTER="global.footer",
	STARTPAGE_CONTENT="startpage.content",
	STARTPAGE_LEFT="startpage.left",
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