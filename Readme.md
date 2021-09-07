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


## Configure your application

```
NAME:
   Micropuzzle - A new cli application

USAGE:
   micropuzzle [global options] command [command options] [arguments...]

DESCRIPTION:
   Application to combine Server Side Include and Afterloading

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --timeoutms value       Timeout for loading Microfrontends (for all slower, it will be use streaming to bring it to the client) (default: 45ms) [$MICROPUZZLE_TIMEOUTMS]
   --logLevel value        Loglevel debug, info, warn, error (default: "info") [$MICROPUZZLE_LOGLEVEL]
   --port value            port where server will be started (default: "3000") [$MICROPUZZLE_PORT]
   --publicfoder value     Folder where all html js css from server directly will be foundable (Public folder for the web) (default: "./public") [$MICROPUZZLE_PUBLICFODER]
   --microfrontends value  A ini file (key=value) key is for logic name of microfrontend. value is the url where to fetch the content (groups are . seperated by using) (default: "./config/frontends.ini") [$MICROPUZZLE_MICROFRONTENDS]
   --fallbackloader value  key of inifile where to find fallbackhtml which will shown if microfrontend is lower than timeout (default: "fallback") [$MICROPUZZLE_FALLBACKLOADER]
   --redisaddr value       The domian/ip:port of redis (default: "localhost:6379") [$MICROPUZZLE_REDISADDR]
   --redisuser value       Username to connect to redis [$MICROPUZZLE_REDISUSER]
   --redispassword value   Password to connect to redis [$MICROPUZZLE_REDISPASSWORD]
   --redisdb value         Db to use by redis (default: 0) [$MICROPUZZLE_REDISDB]
   --managementport value  Port to get data not needed from client (default: 3001) [$MICROPUZZLE_MANAGEMENTPORT]
   --help, -h              show help (default: false)

```


## Configure your possible available Microfrontends
with `--microfrontends` (or over Environment `MICROPUZZLE_MICROFRONTENDS`) you can set the destination of an inifile. 
Inside this file you can configure with logic fragment have which url to load the content. [example](/config/frontends.ini)

## Configure your frontend scarlett
with `--publicfoder` (or over Environment `MICROPUZZLE_PUBLICFODER`) you can set the folder where to find root index.html template
Inside this folder you can use normal html, css, js with should be global there. 

[take a look here](/public/index.html)

## How template html work
You have to set two template flag required: 
- {{.ScriptLoader}} => inside your head Tag to set .js loader
- {{.Loader}} => at the end of the body Tag (important have to be set after each  {{.Reader.Load "..." "..."}})

To load a microfrontend from ini file you have to insert  {{.Reader.Load "key inside your ini file" "logic unique name of this fronend space"}}

The `logic unique name of this fronend space` is importend. About this you have the possibility to change the content of this area. 


## How to change the Content from Client site (example button click)
[look here](/externalServices/footer_old/index.html)

At the end you have to send Custom Javascript events to load new Content
```
  const event = new CustomEvent('load-content',{
    bubbles: true,
    composed: true,
    cancelable: true,
    detail: {
      loading: "key inside your ini file,"
      content: "logic unique name of this fronend space"
  }})
  document.dispatchEvent(event);
```



