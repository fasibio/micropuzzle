package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"gopkg.in/ini.v1"
)

const (
	SocketCommandLoadContent = "LOAD_CONTENT"
	SocketCommandNewContent  = "NEW_CONTENT"
)

type SocketHandler struct {
	cache             ChacheHandler
	Server            *socketio.Server
	proxy             proxy.Proxy
	timeout           time.Duration
	destinations      *ini.File
	fallbackLoaderKey string
}

type LoadContentPayload struct {
	Content     string
	Loading     string
	ExtraHeader http.Header
}

func NewSocketHandler(cache ChacheHandler, timeout time.Duration, destinations *ini.File, fallbackLoaderKey string) SocketHandler {

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	handler := SocketHandler{
		cache:             cache,
		Server:            server,
		proxy:             proxy.Proxy{},
		timeout:           timeout,
		destinations:      destinations,
		fallbackLoaderKey: fallbackLoaderKey,
	}
	server.OnConnect("/", handler.OnConnect)
	server.OnEvent("/", SocketCommandLoadContent, handler.OnLoadContent)
	server.OnError("/", handler.OnError)
	server.OnDisconnect("/", handler.OnDisconnect)
	go func() {
		if err := server.Serve(); err != nil {
			logger.Get().Warnw("socketio listen error", "error", err)
		}
	}()
	return handler
}

type SocketUser struct {
	Connection socketio.Conn
	Id         string
	Roomid     string
}

func (sh *SocketHandler) getSockerUser(s socketio.Conn) SocketUser {
	url := s.URL()
	id := url.Query().Get("streamId")
	return SocketUser{
		Connection: s,
		Id:         s.ID(),
		Roomid:     id,
	}
}
func (sh *SocketHandler) OnConnect(s socketio.Conn) error {
	user := sh.getSockerUser(s)
	sh.Server.JoinRoom("/", user.Roomid, s)
	logs := logger.Get().With("method", "OnConnect", "userID", user.Id, "connectionID", user.Roomid)
	values, err := sh.cache.GetAllValuesForSession(user.Roomid)
	if err != nil {
		logs.Warnw("Error by create connection", "error", err)
		return err
	}
	for _, v := range values {
		sh.Server.BroadcastToRoom("/", user.Roomid, SocketCommandNewContent, NewContentPayload{Key: v.Content, Value: string(v.Value)})
	}
	s.SetContext("")
	logs.Info("New user Connected")
	return nil
}

func (sh *SocketHandler) OnLoadContent(s socketio.Conn, msg LoadContentPayload) {
	user := sh.getSockerUser(s)

	header := s.RemoteHeader()
	for k, v := range msg.ExtraHeader {
		header[k] = v
	}
	result := sh.Load(msg.Loading, msg.Content, user.Roomid, header, s.RemoteAddr().String())
	sh.Server.BroadcastToRoom("/", user.Roomid, SocketCommandNewContent, NewContentPayload{Key: msg.Content, Value: result})
}

func (sh *SocketHandler) OnError(s socketio.Conn, e error) {
	logs := logger.Get().With("method", "OnError")
	logs.Warnw("error at socket connection", "error", e)
}

func (sh *SocketHandler) OnDisconnect(s socketio.Conn, reason string) {
	logs := logger.Get().With("method", "OnDisconnect")
	sh.Server.LeaveAllRooms("/", s)
	logs.Infow("User Disconnect", "error", reason)
}

func (sh *SocketHandler) HandleClientContent(id string, key, value string) {

	if sh.Server.RoomLen("", id) > 0 {
		sh.Server.BroadcastToRoom("/", id, SocketCommandNewContent, NewContentPayload{Key: key, Value: value})
	} else {
		err := sh.cache.Add(id, key, []byte(value))
		if err != nil {
			logs := logger.Get().With("method", "HandleClientContent", "connectionID", id)
			logs.Warnw("error by add data to cache", "error", err)
		}
	}
}

func (sh *SocketHandler) loadAsync(url string, content string, result *chan string, timeout *chan bool, id string, header http.Header, remoteAddr string) {
	res, err := sh.proxy.Get(url, header, remoteAddr)

	if err != nil {
		logger.Get().Warnw("error by load url", "url", url, "error", err)
		return
	}

	contentPage := string(res)
	if len(*timeout) == 1 {
		sh.HandleClientContent(id, content, contentPage)
	} else {
		*result <- contentPage
	}
}

func (sh *SocketHandler) getUrlByFrontendName(name string) string {
	val := strings.Split(name, ".")
	if len(val) == 1 {
		return sh.destinations.Section("").KeysHash()[name]
	}
	return sh.destinations.Section(val[0]).KeysHash()[val[1]]
}

func (sh *SocketHandler) Load(frontend, content string, id string, header http.Header, remoteAddr string) string {
	resultChan := make(chan string, 1)
	timeout := make(chan bool, 1)
	timeoutBubble := make(chan bool, 1)
	go sh.loadAsync(sh.getUrlByFrontendName(frontend), content, &resultChan, &timeoutBubble, id, header, remoteAddr)

	go func() {
		time.Sleep(sh.timeout)
		timeout <- true
	}()
	select {
	case d := <-resultChan:
		{
			return d
		}
	case <-timeout:
		{
			timeoutBubble <- true
			res, err := sh.proxy.Get(sh.getUrlByFrontendName(sh.fallbackLoaderKey), header, remoteAddr)
			if err != nil {
				return "Loading ..."
			}
			return string(res)
		}
	}
}
