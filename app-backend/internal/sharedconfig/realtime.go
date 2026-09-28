package sharedconfig

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// userStreamChannel namespaces this feature's Redis Pub/Sub channels.
func userStreamChannel(username string) string { return "app-backend:user-stream:" + username }

// StartUserStreamRelay starts the one shared Redis subscription every
// registered websocket connection's events flow through, fanning each
// incoming message out to whichever local connections are currently
// registered for that message's username. Call once at boot, after Redis
// is initialized - a no-op if Redis is disabled (matches this feature's
// fire-and-forget design: no socket connection, no live nudge, same as
// today before this existed).
func (gc *GlobalConfig) StartUserStreamRelay() {
	if gc.RedisCache == nil || !gc.RedisCache.Enabled || gc.RedisCache.Client == nil {
		return
	}
	gc.userStreamMutex.Lock()
	if gc.userStreamPubSub != nil {
		gc.userStreamMutex.Unlock()
		return
	}
	// Subscribe with no channels yet - RegisterUserStreamConnection adds
	// per-username channels to this same PubSub object as connections
	// come and go, rather than opening a new Redis subscription per
	// connection.
	pubsub := gc.RedisCache.Client.Subscribe(gc.RedisCache.Context)
	gc.userStreamPubSub = pubsub
	gc.userStreamMutex.Unlock()

	go func() {
		for msg := range pubsub.Channel() {
			username := msg.Channel[len("app-backend:user-stream:"):]

			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("[StartUserStreamRelay] failed to unmarshal message for username [%v]: %v\n", username, err)
				continue
			}

			gc.userStreamMutex.Lock()
			conns := gc.userStreamConnections[username]
			for _, ch := range conns {
				select {
				case ch <- payload:
				default:
					// A slow/backed-up connection must not stall delivery
					// to every other connection - drop for this one
					// reader, same fire-and-forget tradeoff as the rest
					// of this feature.
				}
			}
			gc.userStreamMutex.Unlock()
		}
	}()
}

// RegisterUserStreamConnection registers a new local delivery channel for
// username - fold the returned channel into the websocket handler's
// existing select loop. The first local connection for a username
// subscribes this instance's shared Redis PubSub to that user's channel;
// the last one leaving (via the returned unregister func, call it via
// defer) unsubscribes. Safe to call even when Redis/the relay is
// disabled - the connection just never receives anything over this path
// (push notifications are unaffected).
func RegisterUserStreamConnection(gc *GlobalConfig, username string) (chan map[string]interface{}, func()) {
	ch := make(chan map[string]interface{}, 20)
	if username == "" {
		return ch, func() {}
	}
	connID := uuid.NewString()

	gc.userStreamMutex.Lock()
	if gc.userStreamConnections == nil {
		gc.userStreamConnections = make(map[string]map[string]chan map[string]interface{})
	}
	firstConnectionForUser := len(gc.userStreamConnections[username]) == 0
	if gc.userStreamConnections[username] == nil {
		gc.userStreamConnections[username] = make(map[string]chan map[string]interface{})
	}
	gc.userStreamConnections[username][connID] = ch
	pubsub := gc.userStreamPubSub
	gc.userStreamMutex.Unlock()

	if firstConnectionForUser && pubsub != nil {
		if err := pubsub.Subscribe(gc.RedisCache.Context, userStreamChannel(username)); err != nil {
			log.Printf("[RegisterUserStreamConnection] failed to subscribe for %v: %v\n", username, err)
		}
	}

	unregister := func() {
		gc.userStreamMutex.Lock()
		conns := gc.userStreamConnections[username]
		delete(conns, connID)
		lastConnectionForUser := len(conns) == 0
		if lastConnectionForUser {
			delete(gc.userStreamConnections, username)
		}
		pubsub := gc.userStreamPubSub
		gc.userStreamMutex.Unlock()

		if lastConnectionForUser && pubsub != nil {
			if err := pubsub.Unsubscribe(gc.RedisCache.Context, userStreamChannel(username)); err != nil {
				log.Printf("[RegisterUserStreamConnection] failed to unsubscribe for %v: %v\n", username, err)
			}
		}
	}
	return ch, unregister
}

// PublishUserStreamEvent publishes an event to username's stream. Any
// connection currently registered for username on any instance (via
// RegisterUserStreamConnection, relayed by that instance's
// StartUserStreamRelay) receives it. Fire-and-forget by design, same
// tradeoff as WS3's stream relay: if nobody has a socket open right now,
// the event is simply not delivered here - push notification and durable
// order state are unaffected, this is a live-nudge channel only, not a
// catch-up-able feed.
func (gc *GlobalConfig) PublishUserStreamEvent(username, streamType string, payload interface{}) {
	if username == "" || gc.RedisCache == nil || !gc.RedisCache.Enabled || gc.RedisCache.Client == nil {
		return
	}
	message := map[string]interface{}{"stream": payload, "streamType": streamType}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("[PublishUserStreamEvent] failed to marshal message for %v: %v\n", username, err)
		return
	}
	if err := gc.RedisCache.Client.Publish(gc.RedisCache.Context, userStreamChannel(username), data).Err(); err != nil {
		log.Printf("[PublishUserStreamEvent] failed to publish for %v: %v\n", username, err)
	}
}
