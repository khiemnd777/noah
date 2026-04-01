package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

type JSONHandler func(msg string) error

var (
	pubsubCtx         = context.Background()
	channelHandlers   = make(map[string][]JSONHandler)
	pubsubCancelFuncs = make(map[string]context.CancelFunc)
	messageWgMap      = sync.Map{}
	channelMu         sync.RWMutex
)

type asyncMessage struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func PublishRaw(channel string, message string) error {
	rdb, err := pubsubBackend()
	if err != nil {
		return err
	}
	return rdb.Publish(channel, []byte(message))
}

func PublishJSON(channel string, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return PublishRaw(channel, string(bytes))
}

func PublishJSONAsync(channel string, payload any) error {
	rawBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msgID := uuid.NewString()
	msg := asyncMessage{ID: msgID, Payload: string(rawBytes)}
	wrappedBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	channelMu.RLock()
	handlers := append([]JSONHandler(nil), channelHandlers[channel]...)
	channelMu.RUnlock()

	wg := &sync.WaitGroup{}
	wg.Add(len(handlers))
	messageWgMap.Store(msgID, wg)

	if err := PublishRaw(channel, string(wrappedBytes)); err != nil {
		return err
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
	}

	messageWgMap.Delete(msgID)
	return nil
}

func SubscribeJSON[T any](channel string, handler func(payload *T) error) {
	subscribe(channel, func(msg string) error {
		var result T
		if err := json.Unmarshal([]byte(msg), &result); err != nil {
			return err
		}
		return handler(&result)
	})
}

func SubscribeAsyncJSON[T any](channel string, handler func(payload *T) error) {
	subscribe(channel, func(msg string) error {
		var wrapper asyncMessage
		if err := json.Unmarshal([]byte(msg), &wrapper); err != nil {
			return err
		}

		var result T
		if err := json.Unmarshal([]byte(wrapper.Payload), &result); err != nil {
			return err
		}

		defer func() {
			if value, ok := messageWgMap.Load(wrapper.ID); ok {
				value.(*sync.WaitGroup).Done()
			}
		}()

		return handler(&result)
	})
}

func Unsubscribe(channel string, handlerToRemove JSONHandler) {
	channelMu.Lock()
	defer channelMu.Unlock()

	handlers := channelHandlers[channel]
	newHandlers := make([]JSONHandler, 0, len(handlers))
	for _, h := range handlers {
		if !equalHandlers(h, handlerToRemove) {
			newHandlers = append(newHandlers, h)
		}
	}

	if len(newHandlers) == 0 {
		if cancel, ok := pubsubCancelFuncs[channel]; ok {
			cancel()
			delete(pubsubCancelFuncs, channel)
		}
		delete(channelHandlers, channel)
		return
	}

	channelHandlers[channel] = newHandlers
}

func subscribe(channel string, handler JSONHandler) {
	channelMu.Lock()
	if _, exists := channelHandlers[channel]; !exists {
		channelHandlers[channel] = []JSONHandler{}
	}
	channelHandlers[channel] = append(channelHandlers[channel], handler)
	if _, exists := pubsubCancelFuncs[channel]; exists {
		channelMu.Unlock()
		return
	}
	subCtx, cancel := context.WithCancel(pubsubCtx)
	pubsubCancelFuncs[channel] = cancel
	channelMu.Unlock()

	rdb, err := pubsubBackend()
	if err != nil {
		return
	}

	go func() {
		for {
			sub, err := rdb.Subscribe(subCtx, channel)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			ch := sub.Channel()

			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						_ = sub.Close()
						time.Sleep(time.Second)
						goto retry
					}

					channelMu.RLock()
					handlers := append([]JSONHandler(nil), channelHandlers[channel]...)
					channelMu.RUnlock()

					var wg sync.WaitGroup
					for _, h := range handlers {
						wg.Add(1)
						go func(handler JSONHandler) {
							defer wg.Done()
							_ = handler(string(msg.Payload))
						}(h)
					}
					wg.Wait()
				case <-subCtx.Done():
					_ = sub.Close()
					return
				}
			}
		retry:
		}
	}()
}

func pubsubBackend() (frameworkcache.Backend, error) {
	rdb, err := CacheBackend("pubsub")
	if err == nil {
		return rdb, nil
	}

	rdb, err = CacheBackend("cache")
	if err != nil {
		return nil, fmt.Errorf("redis instance 'pubsub' not available: %w", err)
	}
	return rdb, nil
}

func equalHandlers(a, b JSONHandler) bool {
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}
