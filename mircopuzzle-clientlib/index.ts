export interface PageDeclaration {
  url: string, 
  title: string,
  fragments: {
    [key: string]: string
  }
}

export type PageDeclarations = {
  [key: string]: PageDeclaration
}

export interface HistoryPush{
  url: string | URL;
  title: string;
  state?: any;
}

export const loadMicroFrontend = (area: string, microfrontend: string ) =>{
 
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

export const pushToPage = (page: Page, queryParam?: {[key: string]: string}) =>{
  const p: PageDeclaration = pageDeclarations[page];
  Object.keys(p.fragments).forEach(fragmentKey => {
    loadMicroFrontend(fragmentKey, p.fragments[fragmentKey])
  })
  
  let queryStr = ""
  if (queryParam){
    queryStr = "?"
    queryStr += Object.keys(queryParam).map<string>((key) =>{
      return `${key}=${queryParam[key]}`
    }).join("&")
  }
  

  history.pushState({
    puzzleData: p
  }, p.title,p.url+queryStr)
}