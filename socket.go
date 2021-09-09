package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	promLoadFragmentsTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "micropuzzle_duration_load_milliseconds",
		Help:    "micropuzzle loading nanoseconds for microfrontends",
		Buckets: []float64{1, 5, 10, 30, 50, 80, 100, 1000},
	}, []string{"fragment", "frontend", "afterTimeout", "cached"})
)

func init() {
	prometheus.MustRegister(promLoadFragmentsTime)
}

var upgrader = websocket.Upgrader{} // use default options

type WebSocketHandler struct {
	cache             ChacheHandler
	pubSub            WebSocketBroadcast
	proxy             proxy.Proxy
	timeout           time.Duration
	destinations      Frontends
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
type PubSubNewFragmentPayload struct {
	Payload NewFragmentPayload `json:"payload,omitempty"`
	Id      string             `json:"id,omitempty"`
}

type LoadFragmentPayload struct {
	Content     string `json:"content,omitempty"`
	Loading     string `json:"loading,omitempty"`
	ExtraHeader map[string][]string
}

func (p PubSubNewFragmentPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func getLoggerWithUserInfo(logs *zap.SugaredLogger, user WebSocketUser) *zap.SugaredLogger {
	return logs.With("streamid", user.Id, "address", user.RemoteAddr)
}

func NewWebSocketHandler(cache *RedisHandler, timeout time.Duration, destinations Frontends, fallbackLoaderKey string) WebSocketHandler {
	handler := WebSocketHandler{
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

func (sh *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	user := sh.getSockerUser(c, r)
	logs := getLoggerWithUserInfo(logger.Get(), user)
	sh.user[user.Id] = user
	sh.pubSub.Publish(PubSubCommandNewUser, user.Id)
	if err != nil {
		logs.Warnw("Error by upgrade To websocket con", "error", err)
	}

	go sh.handleMessages(user)

	values, err := sh.cache.GetAllValuesForSession(user.Id)
	if err != nil {
		logs.Warnw("Error by create connection", "error", err)
	}

	for _, v := range values {
		logs.Infow("Send Data to Client found inside cache", "fragment", v.Content)
		sh.updateClientFragment(user.Id, v.Content, string(v.Value))
		go sh.cache.Del(v.Session, v.Content)
	}
}

func (sh *WebSocketHandler) onDelUser(msg *redis.Message, bus WebSocketBroadcast) {
	delete(sh.allKnowUserIds, msg.Payload)
}

func (sh *WebSocketHandler) onNewUser(msg *redis.Message, bus WebSocketBroadcast) {
	sh.allKnowUserIds[msg.Payload] = true
}

func (sh *WebSocketHandler) onNewFragment(msg *redis.Message, bus WebSocketBroadcast) {
	var payload PubSubNewFragmentPayload
	json.Unmarshal([]byte(msg.Payload), &payload)
	user, ok := sh.user[payload.Id]
	if ok {
		err := sh.writeFragmentToClient(user, &payload.Payload)
		if err != nil {
			logger.Get().Warnw("error by send data to client", "error", err, "methode", "onNewFragment")
		}
	}
}

func (sh *WebSocketHandler) getSockerUser(c *websocket.Conn, r *http.Request) WebSocketUser {
	return WebSocketUser{
		Connection:   c,
		Id:           r.URL.Query().Get("streamid"),
		RemoteHeader: r.Header,
		RemoteAddr:   r.RemoteAddr,
	}
}

func (sh *WebSocketHandler) handleMessages(user WebSocketUser) {
	logs := getLoggerWithUserInfo(logger.Get(), user)
	for {
		var messages Message
		err := user.Connection.ReadJSON(&messages)
		if err != nil {
			logs.Infow("error by read json message", "error", err)
			sh.pubSub.Publish(PubSubCommandRemoveUser, user.Id)
			go sh.cache.DelAllForSession(user.Id)
			break
		}
		go sh.interpretMessage(user, messages)
	}
}

func (sh *WebSocketHandler) interpretMessage(user WebSocketUser, msg Message) {
	switch msg.Type {
	case SocketCommandLoadFragment:
		var result LoadFragmentPayload
		mapstructure.Decode(msg.Data, &result)
		sh.onLoadFragment(user, result)
	}
}

func (sh *WebSocketHandler) onLoadFragment(user WebSocketUser, msg LoadFragmentPayload) {

	header := user.RemoteHeader
	for k, v := range msg.ExtraHeader {
		header[k] = v
	}
	result, _, isFallback := sh.LoadFragment(msg.Loading, msg.Content, user.Id, user.RemoteAddr, header)
	sh.writeFragmentToClient(user, &NewFragmentPayload{
		Key:        msg.Content,
		Value:      result,
		IsFallback: isFallback,
	})
}

func (sh *WebSocketHandler) LoadFragmentHandler(w http.ResponseWriter, r *http.Request) {
	fragment := r.URL.Query().Get("fragment")
	frontend := r.URL.Query().Get("frontend")
	userId := r.URL.Query().Get("streamid")
	content, cache, isFallback := sh.LoadFragment(fragment, frontend, userId, r.RemoteAddr, r.Header)
	c, err := json.Marshal(NewFragmentPayload{
		Key:        frontend,
		Value:      content,
		IsFallback: isFallback,
	})
	if err != nil {
		logger.Get().Warnw("error by marshal result", "error", err)
	}
	for k, v := range cache.Header {
		w.Header()[k] = v
	}
	w.Write(c)
}

func (sh *WebSocketHandler) writeFragmentToClient(user WebSocketUser, payload *NewFragmentPayload) error {
	return sh.writeMessage2Client(user, Message{Type: SocketCommandNewContent, Data: payload})
}

func (sh *WebSocketHandler) writeMessage2Client(user WebSocketUser, payload Message) error {
	return user.Connection.WriteJSON(payload)
}

func (sh *WebSocketHandler) updateClientFragment(id, key, value string) {
	_, ok := sh.allKnowUserIds[id]
	if ok {
		err := sh.pubSub.Publish(PubSubCommandNewFragment, PubSubNewFragmentPayload{
			Payload: NewFragmentPayload{
				Key:   key,
				Value: value,
			},
			Id: id})
		if err != nil {
			logger.Get().Warnw("error by publish to redis", "error", err)
		}
	} else {
		err := sh.cache.Add(id, key, value)
		if err != nil {
			logs := logger.Get().With("method", "HandleClientContent", "connectionID", id)
			logs.Warnw("error by add data to cache", "error", err)
		}
	}

}

func (sh *WebSocketHandler) getUrlByFrontendName(name string) string {
	val := strings.Split(name, ".")
	group := "global"
	if len(val) > 1 {
		group = val[0]
	}
	return sh.destinations[group][val[len(val)-1]].Url
}

type AsyncLoadResultChan struct {
	Value string
	Cache proxy.CacheingInformation
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
func (sh *WebSocketHandler) LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheingInformation, bool) {
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
		res, cache, err := sh.proxy.Get(sh.getUrlByFrontendName(sh.fallbackLoaderKey), header, remoteAddr)
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

func (sh *WebSocketHandler) handleFragmentContent(options loadAsyncOptions, content string, cache proxy.CacheingInformation) {
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

func (sh *WebSocketHandler) loadAsync(options loadAsyncOptions) {
	start := time.Now()
	url := sh.getUrlByFrontendName(options.Frontend)
	cachedValue, expire, err := sh.cache.GetPage(options.Frontend)
	fromCache := false
	if err == nil {
		header := make(http.Header)
		header.Set("Cache-Control", fmt.Sprintf("max-age=%.0f;", expire.Seconds()))
		sh.handleFragmentContent(options, cachedValue, proxy.CacheingInformation{Expires: expire, Header: header})
		fromCache = true
	} else {
		res, cache, err := sh.proxy.Get(url, options.Header, options.RemoteAddr)
		if cache.Expires > 0 {
			sh.cache.AddPage(options.Frontend, string(res), cache.Expires)
		}
		if err != nil {
			logger.Get().Warnw("error by load url", "url", url, "error", err)
			return
		}

		fragment := string(res)
		sh.handleFragmentContent(options, fragment, cache)
	}
	sh.writePromMessage(options, fromCache, len(options.Timeout) == 1, start)
}

func (sh *WebSocketHandler) writePromMessage(options loadAsyncOptions, fromCache, insideTimeout bool, start time.Time) {
	promLoadFragmentsTime.WithLabelValues(options.FragmentName, options.Frontend, strconv.FormatBool(insideTimeout), strconv.FormatBool(fromCache)).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
