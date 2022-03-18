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