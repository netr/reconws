package reconws

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_ReadChan(t *testing.T) {
	//chRead := make(chan string)
	//chDone := make(chan bool)
	//cl := NewClient(chRead, chDone, func() {
	//	log.Println("connected")
	//})

	//cl.Connect()
}

var upgrader = websocket.Upgrader{}

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//defer ws.Close()
}
func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()
	s := httptest.NewServer(h)
	wsURL := makeWsProto(s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func TestClient_Close(t *testing.T) {
	ha := handler{}

	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	chRead := make(chan string)
	chDone := make(chan bool)

	cl := NewClient().
		SetChannels(chRead, chDone).
		OnReconnect(func() {
			t.Fatal("should not reconnect")
		})

	ws, err := cl.Connect(makeWsProto(s.URL))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	cl.Close()
	time.Sleep(1 * time.Second)
	close(chRead)
	close(chDone)

	err = ws.UnderlyingConn().Close()
	assert.NotNil(t, err)
}

func TestClient_Reconnecting(t *testing.T) {
	ha := handler{}
	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			t.Fatalf("Close: %v", err)
		}
	}(ws)

	count := 0
	want := 5 // original connection + 4 underlying connection closes in loop
	reconCount := 0
	reconWant := 4 // 4 underlying connection closes in loop
	discCount := 0
	discWant := 5 // 4 underlying connection closes in loop + final cl.Close()

	chRead := make(chan string)
	chDone := make(chan bool)
	cl := NewClient().
		SetChannels(chRead, chDone).
		OnConnect(func() {
			count++
		}).
		OnReconnect(func() {
			reconCount++
		}).
		OnDisconnect(func() {
			discCount++
		})

	_, err := cl.Connect(makeWsProto(s.URL))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	for i := 1; i <= want-1; i++ {
		time.Sleep(100 * time.Millisecond)
		_ = cl.conn.Close()
	}
	time.Sleep(100 * time.Millisecond)
	cl.Close()
	close(chRead)
	close(chDone)

	assert.Equal(t, want, count)
	assert.Equal(t, reconWant, reconCount)
	assert.Equal(t, discWant, discCount)
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
