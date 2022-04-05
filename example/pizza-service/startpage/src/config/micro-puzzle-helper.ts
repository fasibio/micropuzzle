/**
 * Mircopuzzle AUTO-GENERATED CODE: PLEASE DO NOT MODIFY MANUALLY
 */
	
export enum MicropuzzleFrontends {
	ABOUT_CONTENT="about.content",
	CART_CONTENT="cart.content",
	DETAIL_CONTENT="detail.content",
	GLOBAL_FALLBACK="global.fallback",
	GLOBAL_FOOTER="global.footer",
	GLOBAL_HEADER="global.header",
	STARTPAGE_CONTENT="startpage.content",
}


export type Page = 'start'|'cart'|'about'|'detail';
export const pageDeclarations: PageDeclarations = {
  'about': {
    url: '/about.html',
    title: 'about',
    fragments: { 
      'content': 'about.content',
      'footer': 'global.footer',
      'header': 'global.header',  
    }
  },
'cart': {
    url: '/cart.html',
    title: 'Cart',
    fragments: { 
      'content': 'cart.content',
      'footer': 'global.footer',
      'header': 'global.header',  
    }
  },
'detail': {
    url: '/detail.html',
    title: 'detail',
    fragments: { 
      'content': 'detail.content',
      'footer': 'global.footer',
      'header': 'global.header',  
    }
  },
'start': {
    url: '/',
    title: '',
    fragments: { 
      'content': 'startpage.content',
      'footer': 'global.footer',
      'header': 'global.header',  
    }
  },

}
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