package fragments

/**
Subscription methods to handle redis subscriptions
**/

import (
	"encoding/json"

	"github.com/fasibio/micropuzzle/cache"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-redis/redis/v8"
)

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
