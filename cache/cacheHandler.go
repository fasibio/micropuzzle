package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-redis/redis/v8"
)

const MaxTTL = 30 * time.Second

type CacheHandler interface {
	Add(session, key string, value string) error
	Get(session, key string) (DataHolder, error)
	Del(session, key string) error
	DelAllForSession(session string) error
	AddPage(key, value string, ttl time.Duration) error
	GetPage(key string) (string, time.Duration, error)
	AddBlocker(session, key string, value string) error
	GetBlocker(session, key string) (DataHolder, error)
	GetAllValuesForSession(keyPattern string) ([]DataHolder, error)
	DelBlocker(session, key string) error
}

type SubscriptionHandler = func(msg *redis.Message, bus WebSocketBroadcast)

type WebSocketBroadcast interface {
	Publish(channel string, message interface{}) error
	Subscribe() error
	On(channel string, handler SubscriptionHandler)
}

type DataHolder struct {
	Value   string
	Key     string
	Content string
	Session string
}

func (r *RedisHandler) AddPage(key, value string, ttl time.Duration) error {
	res, err := r.zip(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, r.concatPageKey(key), res, ttl).Err()
}

func (r *RedisHandler) GetPage(key string) (string, time.Duration, error) {
	res, err := r.client.Get(r.ctx, r.concatPageKey(key)).Bytes()
	if err != nil {
		return "", time.Duration(0), err
	}
	exp := r.client.TTL(r.ctx, r.concatPageKey(key)).Val()
	val, err := r.unzip(res)
	if err != nil {
		return "", time.Duration(0), err
	}
	return val, exp, nil
}

func (r *RedisHandler) getDataHolderByData(key string, value string) DataHolder {
	keys := strings.Split(strings.Replace(key, fmt.Sprintf("%s_", r.Prefix), "", 1), "_")
	return DataHolder{
		Value:   value,
		Key:     key,
		Session: keys[0],
		Content: keys[1],
	}
}

func (r *RedisHandler) getDataHolderByBlockerData(key string, value string) DataHolder {
	keys := strings.Split(strings.Replace(key, fmt.Sprintf("%s_", r.Prefix), "", 1), "_")
	return DataHolder{
		Value:   value,
		Key:     key,
		Session: keys[1],
		Content: keys[2],
	}
}

type RedisHandler struct {
	client  *redis.Client
	ctx     context.Context
	handler map[string][]SubscriptionHandler
	Prefix  string
}

func NewRedisHandler(opt *redis.Options) (*RedisHandler, error) {
	rdb := redis.NewClient(opt)
	return &RedisHandler{
		client:  rdb,
		ctx:     context.Background(),
		Prefix:  "MICROPUZZLE_",
		handler: make(map[string][]SubscriptionHandler),
	}, nil

}

func (r *RedisHandler) Publish(channel string, message interface{}) error {
	return r.client.Publish(r.ctx, r.PrefixChannel(channel), message).Err()
}

func (r *RedisHandler) PrefixChannel(channel string) string {
	return fmt.Sprintf("%s%s", r.Prefix, channel)
}

func (r *RedisHandler) On(channel string, handler SubscriptionHandler) {
	pChannel := r.PrefixChannel(channel)
	_, ok := r.handler[pChannel]
	if !ok {
		r.handler[pChannel] = []SubscriptionHandler{}
	}
	r.handler[pChannel] = append(r.handler[pChannel], handler)
}

func (r *RedisHandler) Subscribe() error {
	channels := []string{}
	for one := range r.handler {
		channels = append(channels, one)
	}
	res := r.client.Subscribe(r.ctx, channels...)
	subscriptionChan := res.Channel()

	for msg := range subscriptionChan {
		go r.send2Handler(msg)
	}
	return nil
}

func (r *RedisHandler) send2Handler(msg *redis.Message) {
	v, ok := r.handler[msg.Channel]
	if ok {
		for _, one := range v {
			one(msg, r)
		}
	}
}

func (r *RedisHandler) zip(value string) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return []byte{}, err
	}
	if err := gz.Close(); err != nil {
		return []byte{}, err
	}
	return b.Bytes(), nil
}

func (r *RedisHandler) unzip(value []byte) (string, error) {
	var b bytes.Buffer
	b.Write(value)
	gzreader, err := gzip.NewReader(&b)
	if err != nil {
		return "", err
	}
	output, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return "", err
	}
	defer gzreader.Close()

	return string(output), err
}

func (r *RedisHandler) Add(session, key string, value string) error {
	v, err := r.zip(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, r.concatKey(session, key), v, MaxTTL).Err()
}
func (r *RedisHandler) Get(session, key string) (DataHolder, error) {
	res, err := r.client.Get(r.ctx, r.concatKey(session, key)).Bytes()
	if err != nil {
		return DataHolder{}, err
	}
	unZipres, err := r.unzip(res)
	if err != nil {
		return DataHolder{}, err
	}

	return r.getDataHolderByData(r.concatKey(session, key), unZipres), nil
}
func (r *RedisHandler) Del(session, key string) error {
	return r.client.Del(r.ctx, r.concatKey(session, key)).Err()
}

func (r *RedisHandler) DelAllForSession(session string) error {
	keys, err := r.client.Keys(r.ctx, r.concatKey(session, "*")).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.client.Del(r.ctx, keys...).Err()
	}
	return nil
}

func (r *RedisHandler) AddBlocker(session, key string, value string) error {
	return r.client.Set(r.ctx, r.concatBlockerKey(session, key), value, MaxTTL).Err()
}
func (r *RedisHandler) GetBlocker(session, key string) (DataHolder, error) {
	res, err := r.client.Get(r.ctx, r.concatBlockerKey(session, key)).Result()
	if err != nil {
		return DataHolder{}, err
	}
	return r.getDataHolderByBlockerData(r.concatBlockerKey(session, key), res), nil
}

func (r *RedisHandler) DelBlocker(session, key string) error {
	return r.client.Del(r.ctx, r.concatBlockerKey(session, key)).Err()
}

func (r *RedisHandler) GetAllValuesForSession(keyPattern string) ([]DataHolder, error) {
	res, err := r.client.Keys(r.ctx, fmt.Sprintf("%s_%s*", r.Prefix, keyPattern)).Result()
	if err != nil {
		return []DataHolder{}, err
	}
	result := make([]DataHolder, 0)
	for _, one := range res {
		res, err := r.client.Get(r.ctx, one).Bytes()
		if err != nil {
			logger.Get().Warnw("error by get data", "error", err)
		}
		unZipres, err := r.unzip(res)
		if err != nil {
			logger.Get().Warnw("error by unzip data", "error", err)
		}
		result = append(result, r.getDataHolderByData(one, unZipres))
	}
	return result, nil

}
func (r *RedisHandler) concatBlockerKey(session, key string) string {
	return fmt.Sprintf("%s_BLOCKER_%s_%s", r.Prefix, session, key)
}

func (r *RedisHandler) concatPageKey(key string) string {
	return fmt.Sprintf("%s_PAGE_%s", r.Prefix, key)
}

func (r *RedisHandler) concatKey(session, key string) string {
	return fmt.Sprintf("%s_%s_%s", r.Prefix, session, key)
}
