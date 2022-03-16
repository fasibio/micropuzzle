package fragments

/*
	Register function and handler to handle http endpoints

*/

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func (w *FragmentHandler) RegisterHandler(r *chi.Mux, socketPath, socketEndpoint string) {
	r.HandleFunc(fmt.Sprintf("/%s", socketPath), w.handle)
	r.Get(socketEndpoint, w.loadFragmentHandler)
}

func (sh *FragmentHandler) handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	user := sh.getSockerUser(c, r)
	logs := getLoggerWithUserInfo(logger.Get(), user)
	if err != nil {
		logs.Warnw("Error by upgrade To websocket con", "error", err)
	}
	sh.user[user.Id] = user
	err = sh.publishNewUser(user.Id)
	if err != nil {
		logs.Warnw("Error by publish new user", "error", err)
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

func (sh *FragmentHandler) handleMessages(user WebSocketUser) {
	logs := getLoggerWithUserInfo(logger.Get(), user)
	for {
		var messages Message
		err := user.Connection.ReadJSON(&messages)
		if err != nil {
			logs.Debugw("error by read json message", "error", err)
			err := sh.publishRemoveNewUser(user.Id)
			if err != nil {
				logs.Errorw("error by publish remove user", "error", err)
			}
			go sh.cache.DelAllForSession(user.Id)
			break
		}
		go sh.interpretMessage(user, messages)
	}
}

func (sh *FragmentHandler) interpretMessage(user WebSocketUser, msg Message) {
	switch msg.Type {
	case SocketCommandLoadFragment:
		var result LoadFragmentPayload
		mapstructure.Decode(msg.Data, &result)
		sh.onLoadFragment(user, result)
	}
}

func (sh *FragmentHandler) onLoadFragment(user WebSocketUser, msg LoadFragmentPayload) {

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

func (sh *FragmentHandler) loadFragmentHandler(w http.ResponseWriter, r *http.Request) {
	fragment := r.URL.Query().Get("fragment")
	frontend := r.URL.Query().Get("frontend")
	userId := r.Header.Get("streamid")
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

func (sh *FragmentHandler) getSockerUser(c *websocket.Conn, r *http.Request) WebSocketUser {
	return WebSocketUser{
		Connection:   c,
		Id:           r.URL.Query().Get("streamid"),
		RemoteHeader: r.Header,
		RemoteAddr:   r.RemoteAddr,
	}
}
