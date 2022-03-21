package fragments

/**
Publishing methods to push data to redis
**/
import (
	"encoding/json"

	"github.com/fasibio/micropuzzle/cache"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-redis/redis/v8"
)

type PubSubNewFragmentPayload struct {
	Payload newFragmentPayload `json:"payload,omitempty"`
	Id      string             `json:"id,omitempty"`
}

func (p PubSubNewFragmentPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (sh *fragmentHandler) publishNewUser(userid string) error {
	return sh.pubSub.Publish(PubSubCommandNewUser, userid)
}

func (sh *fragmentHandler) publishRemoveNewUser(userid string) error {
	return sh.pubSub.Publish(PubSubCommandRemoveUser, userid)
}

func (sh *fragmentHandler) updateClientFragment(id, key, value string) {
	_, ok := sh.allKnowUserIds[id]
	if ok {
		err := sh.pubSub.Publish(PubSubCommandNewFragment, PubSubNewFragmentPayload{
			Payload: newFragmentPayload{
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

func (sh *fragmentHandler) onDelUser(msg *redis.Message, bus cache.WebSocketBroadcast) {
	delete(sh.allKnowUserIds, msg.Payload)
}

func (sh *fragmentHandler) onNewUser(msg *redis.Message, bus cache.WebSocketBroadcast) {
	sh.allKnowUserIds[msg.Payload] = true
}

func (sh *fragmentHandler) onNewFragment(msg *redis.Message, bus cache.WebSocketBroadcast) {
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
