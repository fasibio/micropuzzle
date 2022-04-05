# Micropuzzle

## It´s a POC at the moment so by careful by using

### A microfrontend Proxy which let you start at minutes to build your Application splittet into fragments.
### Have stable first byte by client always. lazy async load data when content its there

## Quick Getting start. 

Requirements: 
  - Redis

First of all you need some mircofrontend Servers. 

See at folder externalServices the are some examples for different frameworks: 

- [footer](/externalServices/footer) ==> react with styled components


# Complete bad designed example: 
[here](./example/pizza-service/)

You can easy start with `docker-compose up --build` and open your Browser at [localhost:3300](http://localhost:3300).

Also you will find examples there, how local development will look like: 
You can start it with `docker-compose -f docker-compose-dev.yml up --build`.
At [docker-compose-dev.yml](./example/pizza-service/docker-compose-dev.yml) you will see (at line 40) it will use the [frontends-dev.yml](./example/pizza-service/puzzle/config/frontends-dev.yaml)

The Docker Compose will also start a [dockerhost](https://hub.docker.com/r/qoomon/docker-host) where you can call Port at you Host system from a docker container.

At [frontends-dev.yml](./example/pizza-service/puzzle/config/frontends-dev.yaml) one Mircofrontend have the url: `http://dockerhost:3000` for Example [header](./example/pizza-service/header/). Than you have to start header with `yarn dev --host`.

Now you can change the Content of header by Hot-Code reloading inside your Mircopuzzle envoirment. 

This gives you the possibility to develop one Mircofrontend and also see how the other Microfrontends effekt to this one. 


> you will read the `globalOverride: /src` at the header definition also. 
This means, `localhost:3300/src/*` will also forward to this Microfrontend. 
This is often needed to have Hot-Code Reloading


## Configure your application

```
NAME:
   Micropuzzle - A new cli application

USAGE:
   micropuzzle [global options] command [command options] [arguments...]

DESCRIPTION:
   Application to combine Server Side Include and Afterloading

COMMANDS:
   generateType, gen  Generate Typescript types
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --timeoutms value       Timeout for loading Microfrontends (for all slower, it will be use streaming to bring it to the client) (default: 45ms) [$TIMEOUTMS]
   --logLevel value        Loglevel debug, info, warn, error (default: "info") [$LOGLEVEL]
   --port value            port where server will be started (default: "3300") [$PORT]
   --publicfolder value    Folder where all html js css from server directly will be foundable (Public folder for the web) (default: "./public") [$PUBLICFOLDER]
   --microfrontends value  A yaml file to describe available Frontends (default: "./config/frontends.yaml") [$MICROFRONTENDS]
   --fallbackloader value  key of inifile where to find fallbackhtml which will shown if microfrontend is lower than timeout (default: "fallback") [$FALLBACKLOADER]
   --redisaddr value       The domian/ip:port of redis (default: "localhost:6379") [$REDISADDR]
   --redisuser value       Username to connect to redis [$REDISUSER]
   --redispassword value   Password to connect to redis [$REDISPASSWORD]
   --redisdb value         Db to use by redis (default: 0) [$REDISDB]
   --managementport value  Port to get data not needed from client (default: 3301) [$MANAGEMENTPORT]
   --help, -h              show help (default: false)

```


## Configure your possible available Microfrontends
with `--microfrontends` (or over Environment `MICROFRONTENDS`) you can set the destination of an yaml-file. 
Inside this file you can configure with logic fragment have which url to load the content. [example](/config/frontends.yaml)

## Configure your frontend skeleton
with `--publicfolder` (or over Environment `PUBLICFOLDER`) you can set the folder where to find root html templates
Inside this folder you can use normal html, css, js with should be global there. 

[take a look here](/public/index.html)

**Description:**

```yaml
version: 1 # The version of configuration (at moment there is only version 1)


definitions: ## All declaration of Microfrontends 
  global: ## group name
    footer: ## frontend name (the logic name will be than groupName.FrontendName so here: global.footer)
      url: http://localhost:5000 # The Url to Mircofrontend (From Server side ==> this url will be called from Micropuzzle server)
    fallback:
      url: http://localhost:3300/loader.html
    header: 
      url: http://header
  startpage:
    content:
      url: http://localhost:5001
    left:
      url: http://localhost:5003
  about:
    content:
      url: http://localhost:5002

pages: # All page declarations
  global: # reserved pagename 
    fragments: # Here you can define fragments which will be added, as default, to each page if it is not defined there
      # Key is the name you habe defined inside your html template ( {{.Reader.Load "header"}}). 
              # Value is the name you defined inside you definitions object. (See 13 Lines up)
      header: global.header 
      footer: global.footer
    template: ./index.html # Also the html default html template to use. If you not specifier one special for a page this will be used
  start: # A Page
    url: "/" # The Url of the page to identifier
    title: "" # title of the page (will be ignored at the moment from most browser)
    fragments: # All page spezific Mircofrontends, it will be extends with the global header and footer
      content: startpage.content
  about: 
    url: "/about.html"
    title: "about"
    fragments: 
      content: about.content

```

## How template html work
You have to set two template flag required: 
- {{.ScriptLoader}} => inside your head Tag to set .js loader
- {{.Loader}} => at the end of the body Tag (important have to be set after each  {{.Reader.Load "logic unique name of this fronend space"}})

The `logic unique name of this fronend space` is importend. About this you have the possibility to change the content of this area. 

Inside your Microfrontends yaml there is a group called `pages`. 
Here you can define which  page(url) have to load which microfrontend for `the logic unique area`


## How to change the Content from Client site (example button click)
[look here](/externalServices/footer_old/index.html)

At the end you have to send Custom Javascript events to load new Content. 
But a Typescript snippet will be generated for you. 
Unter the Managmenturl (default 3301). 
So if you start it local: 
http://localhost:3301/micro-puzzle-helper.ts

This snippet you can include into your Microfrontends. [example](./example/pizza-service/header/src/App.tsx): 
- See generation Script (`genTypes`) at [package.json](./example/pizza-service/header/package.json)

This Code snippet will also have enum MicropuzzleFrontends and pages where you find all keys from your [frontends.yaml](/config/frontends.yaml)
