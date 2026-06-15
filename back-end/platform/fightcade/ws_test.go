package fightcade

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// startTestServer creates a WebSocket server that calls handler for each message.
// Returns the wsClient connected to it.
func startTestServer(t *testing.T, handler func(conn *websocket.Conn, msg map[string]any)) *wsClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade: %v", err)
		}
		defer conn.Close()
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var msg map[string]any
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}
			handler(conn, msg)
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	c := &wsClient{conn: conn, done: make(chan struct{})}
	go c.recvLoop()
	t.Cleanup(func() { c.close() })
	return c
}

func reply(conn *websocket.Conn, msg map[string]any) {
	data, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, data)
}

func TestSendCmd_CorrelatesResponse(t *testing.T) {
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
			"echo":       msg["req"],
		})
	})

	ctx := context.Background()
	resp, err := client.sendCmd(ctx, map[string]any{"req": "ping"})
	if err != nil {
		t.Fatalf("sendCmd: %v", err)
	}
	if resp["echo"] != "ping" {
		t.Errorf("expected echo=ping, got %v", resp["echo"])
	}
}

func TestSendCmd_MultipleInFlight(t *testing.T) {
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		idx := msg["requestIdx"]
		if msg["req"] == "slow" {
			time.Sleep(50 * time.Millisecond)
		}
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": idx,
			"req":        msg["req"],
		})
	})

	ctx := context.Background()
	var wg sync.WaitGroup
	results := make([]map[string]any, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		resp, _ := client.sendCmd(ctx, map[string]any{"req": "slow"})
		results[0] = resp
	}()
	go func() {
		defer wg.Done()
		resp, _ := client.sendCmd(ctx, map[string]any{"req": "fast"})
		results[1] = resp
	}()
	wg.Wait()

	if results[0]["req"] != "slow" {
		t.Errorf("slow request got wrong response: %v", results[0]["req"])
	}
	if results[1]["req"] != "fast" {
		t.Errorf("fast request got wrong response: %v", results[1]["req"])
	}
}

func TestSendCmd_Timeout(t *testing.T) {
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		// Never reply.
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.sendCmd(ctx, map[string]any{"req": "hang"})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestSendFire_SetsRequestIdxNegOne(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
	})

	err := client.sendFire(map[string]any{"req": "leave", "channelname": "test"})
	if err != nil {
		t.Fatalf("sendFire: %v", err)
	}

	select {
	case msg := <-received:
		if idx, ok := msg["requestIdx"].(float64); !ok || idx != -1 {
			t.Errorf("expected requestIdx=-1, got %v", msg["requestIdx"])
		}
		if msg["req"] != "leave" {
			t.Errorf("expected req=leave, got %v", msg["req"])
		}
	case <-time.After(time.Second):
		t.Fatal("server never received message")
	}
}

func TestOn_DispatchesPushedEvents(t *testing.T) {
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		// On any request, push back an event (no result field).
		reply(conn, map[string]any{
			"req":  "challenge",
			"user": map[string]any{"name": "opponent"},
		})
		// Then send the actual response.
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
		})
	})

	got := make(chan string, 1)
	client.on("challenge", func(msg map[string]any) {
		user, _ := msg["user"].(map[string]any)
		name, _ := user["name"].(string)
		got <- name
	})

	ctx := context.Background()
	_, _ = client.sendCmd(ctx, map[string]any{"req": "join"})

	select {
	case name := <-got:
		if name != "opponent" {
			t.Errorf("expected opponent, got %s", name)
		}
	case <-time.After(time.Second):
		t.Fatal("event handler never called")
	}
}

func TestDone_ClosesOnServerDisconnect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	c := &wsClient{conn: conn, done: make(chan struct{})}
	go c.recvLoop()

	select {
	case <-c.done:
	case <-time.After(time.Second):
		t.Fatal("done channel never closed after server disconnect")
	}
}

func TestLogin_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
			"cookie":     "abc123",
		})
	})

	ctx := context.Background()
	resp, err := client.login(ctx, "player1", "secret")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp["cookie"] != "abc123" {
		t.Errorf("expected cookie abc123, got %v", resp["cookie"])
	}

	msg := <-received
	if msg["req"] != "login" {
		t.Errorf("expected req=login, got %v", msg["req"])
	}
	if msg["username"] != "player1" {
		t.Errorf("expected username=player1, got %v", msg["username"])
	}
	if msg["userpass"] != "secret" {
		t.Errorf("expected userpass=secret, got %v", msg["userpass"])
	}
}

func TestAutoLogin_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
		})
	})

	ctx := context.Background()
	_, err := client.autoLogin(ctx, "mycookie")
	if err != nil {
		t.Fatalf("autoLogin: %v", err)
	}

	msg := <-received
	if msg["req"] != "autologin" {
		t.Errorf("expected req=autologin, got %v", msg["req"])
	}
	if msg["cookie"] != "mycookie" {
		t.Errorf("expected cookie=mycookie, got %v", msg["cookie"])
	}
}

func TestJoinChannel_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
			"users":      []any{},
		})
	})

	ctx := context.Background()
	_, err := client.joinChannel(ctx, "sfiii3nr1")
	if err != nil {
		t.Fatalf("joinChannel: %v", err)
	}

	msg := <-received
	if msg["req"] != "join" {
		t.Errorf("expected req=join, got %v", msg["req"])
	}
	if msg["channelname"] != "sfiii3nr1" {
		t.Errorf("expected channelname=sfiii3nr1, got %v", msg["channelname"])
	}
	if msg["status"] != "available" {
		t.Errorf("expected status=available, got %v", msg["status"])
	}
}

func TestLeaveChannel_IsFireAndForget(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
	})

	err := client.leaveChannel("sfiii3nr1")
	if err != nil {
		t.Fatalf("leaveChannel: %v", err)
	}

	select {
	case msg := <-received:
		if msg["req"] != "leave" {
			t.Errorf("expected req=leave, got %v", msg["req"])
		}
		if idx, _ := msg["requestIdx"].(float64); idx != -1 {
			t.Errorf("expected requestIdx=-1, got %v", idx)
		}
	case <-time.After(time.Second):
		t.Fatal("server never received leave")
	}
}

func TestSearchChannels_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
			"channels":   []any{},
		})
	})

	ctx := context.Background()
	_, err := client.searchChannels(ctx, "street fighter", 2)
	if err != nil {
		t.Fatalf("searchChannels: %v", err)
	}

	msg := <-received
	if msg["req"] != "channels" {
		t.Errorf("expected req=channels, got %v", msg["req"])
	}
	if msg["filter"] != "street fighter" {
		t.Errorf("expected filter='street fighter', got %v", msg["filter"])
	}
	if page, _ := msg["page"].(float64); page != 2 {
		t.Errorf("expected page=2, got %v", page)
	}
}

func TestChallengeUser_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
		})
	})

	ctx := context.Background()
	_, err := client.challengeUser(ctx, "opponent", "sfiii3nr1", 3, 5)
	if err != nil {
		t.Fatalf("challengeUser: %v", err)
	}

	msg := <-received
	if msg["req"] != "challenge" {
		t.Errorf("expected req=challenge, got %v", msg["req"])
	}
	if msg["username"] != "opponent" {
		t.Errorf("expected username=opponent, got %v", msg["username"])
	}
	if cid, _ := msg["challengeid"].(float64); cid != 3 {
		t.Errorf("expected challengeid=3, got %v", cid)
	}
	if ranked, _ := msg["ranked"].(float64); ranked != 5 {
		t.Errorf("expected ranked=5 (FT5), got %v", ranked)
	}
}

func TestAcceptChallenge_SendsCorrectPayload(t *testing.T) {
	received := make(chan map[string]any, 1)
	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		received <- msg
		reply(conn, map[string]any{
			"result":     float64(200),
			"requestIdx": msg["requestIdx"],
		})
	})

	ctx := context.Background()
	_, err := client.acceptChallenge(ctx, "challenger", "sfiii3nr1", 5, 3)
	if err != nil {
		t.Fatalf("acceptChallenge: %v", err)
	}

	msg := <-received
	if msg["req"] != "accept" {
		t.Errorf("expected req=accept, got %v", msg["req"])
	}
	if msg["username"] != "challenger" {
		t.Errorf("expected username=challenger, got %v", msg["username"])
	}
	if cid, _ := msg["challengeid"].(float64); cid != 5 {
		t.Errorf("expected challengeid=5, got %v", cid)
	}
	if ranked, _ := msg["ranked"].(float64); ranked != 3 {
		t.Errorf("expected ranked=3 (FT3), got %v", ranked)
	}
}
