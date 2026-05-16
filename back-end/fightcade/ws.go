package fightcade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	serverURI  = "wss://ggs.fightcade.com/ws/"
	userAgent  = "Fightcade2-OSX-v2.1.45"
	reqTimeout = 15 * time.Second
)

type eventHandler func(msg map[string]any)

type wsClient struct {
	conn       *websocket.Conn
	requestIdx atomic.Int64
	pending    sync.Map // int64 -> chan map[string]any
	handlers   sync.Map // string -> eventHandler
	done       chan struct{}
}

func connect(ctx context.Context) (*wsClient, error) {
	dialer := websocket.Dialer{}
	header := http.Header{"User-Agent": {userAgent}}
	conn, _, err := dialer.DialContext(ctx, serverURI, header)
	if err != nil {
		return nil, fmt.Errorf("websocket dial: %w", err)
	}
	c := &wsClient{conn: conn, done: make(chan struct{})}
	go c.recvLoop()
	return c, nil
}

func (c *wsClient) close() {
	c.conn.Close()
	<-c.done
}

func (c *wsClient) on(event string, handler eventHandler) {
	c.handlers.Store(event, handler)
}

func (c *wsClient) sendCmd(ctx context.Context, payload map[string]any) (map[string]any, error) {
	idx := c.requestIdx.Add(1) - 1
	payload["requestIdx"] = idx
	ch := make(chan map[string]any, 1)
	c.pending.Store(idx, ch)
	defer c.pending.Delete(idx)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *wsClient) sendFire(payload map[string]any) error {
	payload["requestIdx"] = -1
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *wsClient) recvLoop() {
	defer close(c.done)
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var msg map[string]any
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if _, hasResult := msg["result"]; hasResult {
			if idxf, ok := msg["requestIdx"]; ok {
				idx := int64(idxf.(float64))
				if val, loaded := c.pending.LoadAndDelete(idx); loaded {
					val.(chan map[string]any) <- msg
				}
			}
		} else {
			req, _ := msg["req"].(string)
			if val, ok := c.handlers.Load(req); ok {
				val.(eventHandler)(msg)
			}
		}
	}
}

func (c *wsClient) login(ctx context.Context, username, password string) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":      "login",
		"username": username,
		"userpass": password,
		"location": map[string]any{},
	})
}

func (c *wsClient) autoLogin(ctx context.Context, cookie string) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":      "autologin",
		"cookie":   cookie,
		"location": map[string]any{},
	})
}

func (c *wsClient) joinChannel(ctx context.Context, channelname string) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":         "join",
		"channelname": channelname,
		"status":      "available",
		"away":        false,
		"idx":         -1,
	})
}

func (c *wsClient) leaveChannel(channelname string) error {
	return c.sendFire(map[string]any{
		"req":         "leave",
		"channelname": channelname,
	})
}

func (c *wsClient) searchChannels(ctx context.Context, query string, page int) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":       "channels",
		"filter":    query,
		"paginated": true,
		"page":      page,
	})
}

func (c *wsClient) challengeUser(ctx context.Context, username, channelname string, challengeid int, ranked bool) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":         "challenge",
		"username":    username,
		"channelname": channelname,
		"challengeid": challengeid,
		"ranked":      ranked,
	})
}

func (c *wsClient) acceptChallenge(ctx context.Context, username, channelname string, challengeid int, ranked bool) (map[string]any, error) {
	return c.sendCmd(ctx, map[string]any{
		"req":         "accept",
		"username":    username,
		"channelname": channelname,
		"challengeid": challengeid,
		"ranked":      ranked,
	})
}
