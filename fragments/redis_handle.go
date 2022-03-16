package fragments

/**
Publishing methods to push data to redis
**/
import (
	"encoding/json"

	"github.com/fasibio/micropuzzle/logger"
)

type PubSubNewFragmentPayload struct {
	Payload NewFragmentPayload `json:"payload,omitempty"`
	Id      string             `json:"id,omitempty"`
}

func (p PubSubNewFragmentPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (sh *FragmentHandler) publishNewUser(userid string) error {
	return sh.pubSub.Publish(PubSubCommandNewUser, userid)
}

func (sh *FragmentHandler) publishRemoveNewUser(userid string) error {
	return sh.pubSub.Publish(PubSubCommandRemoveUser, userid)
}

func (sh *FragmentHandler) updateClientFragment(id, key, value string) {
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
