package pubsub

import (
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type JsonHandler = frameworkruntime.JSONHandler

func publish(channel string, message string) error {
	return frameworkruntime.PublishRaw(channel, message)
}

func Publish(channel string, payload any) error {
	return frameworkruntime.PublishJSON(channel, payload)
}

func PublishAsync(channel string, payload any) error {
	return frameworkruntime.PublishJSONAsync(channel, payload)
}

func Subscribe[T any](channel string, handler func(payload *T) error) {
	frameworkruntime.SubscribeJSON(channel, handler)
}

func SubscribeAsync[T any](channel string, handler func(payload *T) error) {
	frameworkruntime.SubscribeAsyncJSON(channel, handler)
}

func Unsubscribe(channel string, handlerToRemove JsonHandler) {
	frameworkruntime.Unsubscribe(channel, handlerToRemove)
}
