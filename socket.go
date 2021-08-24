package main

import (
	"github.com/fasibio/micropuzzle/logger"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type SocketHandler struct {
	cache  ChacheHandler
	Server *socketio.Server
}

type LoadContentPayload struct {
	Content string
	Loading string
}

func NewSocketHandler(cache ChacheHandler) SocketHandler {

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
		cache:  cache,
		Server: server,
	}
	server.OnConnect("/", handler.OnConnect)
	server.OnEvent("/", "LOAD_CONTENT", handler.OnLoadContent)
	server.OnError("/", handler.OnError)
	server.OnDisconnect("/", handler.OnDisconnect)
	go func() {
		if err := server.Serve(); err != nil {
			logger.Get().Warnw("socketio listen error", "error", err)
		}
	}()
	return handler
}

func (sh *SocketHandler) OnConnect(s socketio.Conn) error {
	url := s.URL()
	id := url.Query().Get("streamId")
	sh.Server.JoinRoom("/", id, s)
	logs := logger.Get().With("method", "OnConnect", "userID", s.ID(), "connectionID", id)
	values, err := sh.cache.GetAllValuesForSession(id)
	if err != nil {
		logs.Warnw("Error by create connection", "error", err)
		return err
	}
	for _, v := range values {
		sh.Server.BroadcastToRoom("/", id, "NEW_CONTENT", NewContentPayload{Key: v.Content, Value: string(v.Value)})
	}
	s.SetContext("")
	logs.Info("New user Connected")
	return nil
}

func (sh *SocketHandler) OnLoadContent(s socketio.Conn, msg LoadContentPayload) {
	//@TODO IMPLEMENT
	logs := logger.Get().With("method", "OnLoadContent")
	logs.Warn("LoadContent NOT IMPLEMENT NOW")
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
