@img/logo.png 

## Today we are talking about:

- How do I define microfrontends
- What are they needed for?
- The customer/visitor should not pay for microfrontends!
- Keep the development as easy as possible
- Overview of “Micropuzzle
- Micropuzzle Routing 

## We are not talking about: 

- make microfrontends sense at general

@img/im-sorry-does-the-cute-cat-help.jpg

# How do I define microfrontends

 Similar to microservices, a microfrontend system should have:

- No deployment dependencies 

- A clean context separation

- A defined communication way between separations

## whats not a part of my mircofrontend definition

- Reuse of some components at different points

# What are they needed for?

- companies where more than one team work at the same frontend. 

 ==> and have for example problem at deployment process 

- Very huge webpages which should seperate into different contexts to have less cross impacts

# the customer/visitor should not pay for microfrontends!

we should create a solution where the customer does not: 

- Have a much larger download size as monfrontends

- The same possibilities to find the page (SEO)

 Most of them is dependant to business logic

# keep the development as easy as possible

I think nothing more to say ...

# conclusion

For me microfrontends and Server Side Include (SSI) go hand in hand

# Back to

@img/logo.png 

# How Micropuzzle will be used: 

It´s a dockercontainer... 

## Some important things to know:

## 1st

You always need a Redis to start Micropuzzle.
 For Caching and to share information between replications.

## 2nd

You need to add a public folder. (All global HTML stuff is inside and root index.html as well)

## 3rd

You need to add a configuration YAML with the target configuration

# A configuration example

# How the example looks like:

@img/micropuzzle_presentationPic3.png

## 1st YAML File

@code/yml
  global:
    header:
      url: http://header.interal
    fallback:
      url: http://micropuzzleUrl/loader.html
    footer: 
      url: http://footer.interal
    sidebar:
      url: http://sidebar.interal
  start:
    content:
      url: http://startContent.interal
  detail:
    content:
      url: http://detailContent.interal

this will generate useable links: 

@code/txt
  - header
  - fallback
  - footer
  - sidebar
  - start.content
  - detail.content

## 2nd html code

index.html
@code/html
  <head> 
     ...
     <!-- Load needed Micropuzzle specific 
          Client Javascript Libs
     -->
     {{.ScriptLoader}} 
  </head>
  <body>
    <header>
      <!-- 
          Add the result of http://header.interal here 
          and give this area the logic name header-area
      -->
      {{.Reader.Load "header" "header-area"}}
    </header>
    <main class="left">
      <!-- 
          Add the result of http://sidebar.interal here 
          and give this area the logic name left
      -->
      {{.Reader.Load "sidebar" "left"}}
    </main>
    <main class="content">
      <!-- 
          Add the result of http://startContent.interal here 
          and give this area the logic name content
      -->
      {{.Reader.Load "start.content" "content"}}
    </main>
    <footer>
      <!-- 
        Add the result of http://footer.interal here
        and give this area the logic name footer
      -->
      {{.Reader.Load "footer" "footer"}}
    </footer>
    <!-- 
      Has to be set after all .Reader.Load calls 
      will include some micropuzzle spezific html code 
      and will start loading on server side
    -->
    {{.Loader}}
  </body>

this file will be a part of the public folder

to complete this example: 

loader.html
@code/html
  <style>
    .skeleton-box {
      display: inline-block;
      height: 1em;
      position: relative;
      overflow: hidden;
      background-color: #DDDBDD;
    }
    @keyframes shimmer {
      100% {
        transform: translateX(100%);
      }
    }
    .skeleton-box::after {
      position: absolute;
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      transform: translateX(-100%);
      background-image: linear-gradient(
        90deg, 
      rgba(255, 255, 255, 0) 0,
      rgba(255, 255, 255, 0.2) 20%,
      rgba(255, 255, 255, 0.5) 60%,
      rgba(255, 255, 255, 0));
      animation: shimmer 2s infinite;
      content: '';
    }
  </style>
  <div class="skeleton-box" style="width: 300px;height: 500px;"></div>

this file will be also a part of the public folder

# Routing 

To navigate from *start.content* to *about.content* some JavaScript libs are needed

It works with custom events. At the end the client has to execute a CustomEvent on its side.  

@code/javascript
  const event = new CustomEvent('load-content',{
    bubbles: true,
    composed: true,
    cancelable: true,
    detail: {
      // the name of the Microfronted (result of yml)
      loading: 'about.content', 
      // the logic named area where to exchange
      content: 'content' 
    }
  })
  document.dispatchEvent(event);

# Example of the Flow

Best case would be SSI + Router

@img/micropuzzle_happy_path.png

But we will have an unhappy path as well…

@img/micropuzzle_presentationPic1.png

@img/micropuzzle_presentationPic2.png

problem is solved with Micropuzzle

@img/micropuzzle_unhappy_path_1.png

# Routing 

For Routing exactly the same happens

if the microfrontend takes too long, the client will open a Websocket Connection

