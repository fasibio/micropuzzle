package fragments

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fasibio/micropuzzle/cache"
	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/fasibio/micropuzzle/resultmanipulation"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type FragmentHandler struct {
	cache             cache.CacheHandler
	pubSub            cache.WebSocketBroadcast
	proxy             proxy.Proxy
	timeout           time.Duration
	destinations      configloader.Frontends
	fallbackLoaderKey string
	user              map[string]WebSocketUser
	allKnowUserIds    map[string]bool
}

type Message struct {
	Type string      `json:"type,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type WebSocketUser struct {
	Connection   *websocket.Conn
	Id           string
	RemoteHeader http.Header
	RemoteAddr   string
}

type loadAsyncOptions struct {
	Frontend     string
	FragmentName string
	UserId       string
	Uuid         string
	RemoteAddr   string
	Header       http.Header
	Result       chan<- AsyncLoadResultChan
	Timeout      chan<- bool
}

type AsyncLoadResultChan struct {
	Value string
	Cache proxy.CacheingInformation
}

const (
	SocketCommandLoadFragment = "LOAD_CONTENT"
	SocketCommandNewContent   = "NEW_CONTENT"
	PubSubCommandNewFragment  = "new_fragment" //PubSubNewFragmentPayload
	PubSubCommandNewUser      = "new_user"     // string ==> streamId
	PubSubCommandRemoveUser   = "remove_user"  //string ==> streamId
)

type NewFragmentPayload struct {
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	IsFallback bool   `json:"isFallback,omitempty"`
}

type LoadFragmentPayload struct {
	Content     string `json:"content,omitempty"`
	Loading     string `json:"loading,omitempty"`
	ExtraHeader map[string][]string
}

func NewFragmentHandler(cache *cache.RedisHandler, timeout time.Duration, destinations configloader.Frontends, fallbackLoaderKey string) FragmentHandler {
	handler := FragmentHandler{
		cache:             cache,
		pubSub:            cache,
		proxy:             proxy.Proxy{},
		timeout:           timeout,
		destinations:      destinations,
		fallbackLoaderKey: fallbackLoaderKey,
		user:              make(map[string]WebSocketUser),
		allKnowUserIds:    make(map[string]bool),
	}

	handler.pubSub.On(PubSubCommandNewFragment, handler.onNewFragment)
	handler.pubSub.On(PubSubCommandNewUser, handler.onNewUser)
	handler.pubSub.On(PubSubCommandRemoveUser, handler.onDelUser)
	go handler.pubSub.Subscribe()
	return handler
}

func (sh *FragmentHandler) writeFragmentToClient(user WebSocketUser, payload *NewFragmentPayload) error {
	return user.Connection.WriteJSON(Message{Type: SocketCommandNewContent, Data: payload})
}

// LoadFragment try to load the microfrontend
// it this need longer than defined timeout it will return a fallback(some loader) instance of microfrontend.
// It will be start an asnyc loader to get data over the websocket connection to the client if it is there
// frontend the key defined at destinations.
// fragmentName the part of conent to load this
// userId the uuid from client
// remoteAddr which comes from client (needed for proxy)
// header comes from client
// It retruns the content and a bool if is a fallback and not the microfrontent content
func (sh *FragmentHandler) LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheingInformation, bool) {
	uuid, err := uuid.NewV4()
	if err != nil {
		logger.Get().Warnw("Unexpected error happens by gernerate uuid", "error", err)
	}
	err = sh.cache.AddBlocker(userId, fragmentName, uuid.String())
	if err != nil {
		logger.Get().Warnw("Error by write blocker", "error", err)
	}
	resultChan := make(chan AsyncLoadResultChan, 1)
	timeout := make(chan bool, 1)
	timeoutBubble := make(chan bool, 1)
	go sh.loadAsync(loadAsyncOptions{
		Frontend:     frontend,
		FragmentName: fragmentName,
		UserId:       userId,
		RemoteAddr:   remoteAddr,
		Uuid:         uuid.String(),
		Header:       header,
		Result:       resultChan,
		Timeout:      timeoutBubble,
	})

	go func() {
		time.Sleep(sh.timeout)
		timeout <- true
	}()
	select {
	case d := <-resultChan:
		return d.Value, d.Cache, false
	case <-timeout:
		start := time.Now()
		timeoutBubble <- true
		cachedValue, _, err := sh.cache.GetPage(sh.fallbackLoaderKey)
		if err == nil {
			sh.writePromMessage(loadAsyncOptions{
				FragmentName: fragmentName,
				Frontend:     sh.fallbackLoaderKey,
			}, true, true, start)
			promLoadFragmentsTime.WithLabelValues(fragmentName, sh.fallbackLoaderKey, "true", "true").Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
			return cachedValue, proxy.CacheingInformation{Expires: time.Duration(0)}, true
		}
		res, cache, err := sh.proxy.Get(sh.destinations.GetUrlByFrontendName(sh.fallbackLoaderKey), header, remoteAddr)
		if cache.Expires > 0 {
			sh.cache.AddPage(sh.fallbackLoaderKey, string(res), cache.Expires)
		}
		sh.writePromMessage(loadAsyncOptions{
			FragmentName: fragmentName,
			Frontend:     sh.fallbackLoaderKey,
		}, false, true, start)
		if err != nil {
			return "Loading ...", proxy.CacheingInformation{Expires: time.Duration(0)}, true
		}

		return string(res), proxy.CacheingInformation{Expires: time.Duration(0)}, true
	}
}

func (sh *FragmentHandler) handleFragmentContent(options loadAsyncOptions, content string, cache proxy.CacheingInformation) {
	data, err := sh.cache.GetBlocker(options.UserId, options.FragmentName)
	if err != nil {
		logger.Get().Infow("Error by get blockerdata from cache this is not an error at all it also could mean other content was faster at loading", "error", err, "FragmentName", options.FragmentName)
	}
	if string(data.Value) == options.Uuid {
		if len(options.Timeout) == 1 {
			sh.updateClientFragment(options.UserId, options.FragmentName, content)
		} else {
			options.Result <- AsyncLoadResultChan{
				Value: content,
				Cache: cache,
			}
		}
	}
}

func (sh *FragmentHandler) loadAsync(options loadAsyncOptions) {
	start := time.Now()
	url := sh.destinations.GetUrlByFrontendName(options.Frontend)
	cachedValue, expire, err := sh.cache.GetPage(options.Frontend)
	fromCache := false
	if err == nil {
		header := make(http.Header)
		header.Set("Cache-Control", fmt.Sprintf("max-age=%.0f;", expire.Seconds()))
		sh.handleFragmentContent(options, cachedValue, proxy.CacheingInformation{Expires: expire, Header: header})
		fromCache = true
	} else {
		res, cache, err := sh.proxy.Get(url, options.Header, options.RemoteAddr)
		if err != nil {
			logger.Get().Warnw("Error by get data from Microserviceurl", "error", err)
			return
		}
		data, err := resultmanipulation.ChangePathOfRessources(string(res), options.Frontend)
		if err != nil {
			logger.Get().Warnw("Error by change path of ressources", "error", err)
			data = (string(res))
		}
		if cache.Expires > 0 {
			sh.cache.AddPage(options.Frontend, data, cache.Expires)
		}
		if err != nil {
			logger.Get().Warnw("error by load url", "url", url, "error", err)
			return
		}

		fragment := string(data)
		sh.handleFragmentContent(options, fragment, cache)
	}
	sh.writePromMessage(options, fromCache, len(options.Timeout) == 1, start)
}
