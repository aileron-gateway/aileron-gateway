//go:build integration
// +build integration

package httpproxy_test

import (
	"crypto/md5"
	"crypto/rand"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"time"

	"golang.org/x/net/websocket"
)

type wrappedResponseWriter struct {
	*httptest.ResponseRecorder
	bodies     []string
	times      []time.Time
	milliTimes []int64
}

func (w *wrappedResponseWriter) Flush() {
	w.ResponseRecorder.Flush()
	b := make([]byte, 100)
	n, _ := w.ResponseRecorder.Body.Read(b)
	w.bodies = append(w.bodies, string(b[:n]))

	// Record times.
	now := time.Now()
	w.times = append(w.times, now)
	w.milliTimes = append(w.milliTimes, now.UnixMilli())
}

type fixedBodyHandler struct {
	body string
}

func (h *fixedBodyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.body))
}

// chunkedResponseHandler is a http handler that returns a chunked response.
type returnChunkedHandler struct {
	header   http.Header
	interval time.Duration
	bodies   []string
}

func (h *returnChunkedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	maps.Copy(w.Header(), h.header)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()
	for _, b := range h.bodies {
		w.Write([]byte(b))
		flusher.Flush()
		time.Sleep(h.interval)
	}
}

// chunkedResponseHandler is a http handler that receives chunked requests.
type receiveChunkedHandler struct {
	bodies     []string
	times      []time.Time
	milliTimes []int64
}

func (h *receiveChunkedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 100)
	for {
		n, err := r.Body.Read(b)
		if n > 0 {
			h.bodies = append(h.bodies, string(b[:n]))
			now := time.Now()
			h.times = append(h.times, now)
			h.milliTimes = append(h.milliTimes, now.UnixMilli())
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}
	w.Write([]byte("ok"))
}

// webSocketHandler receives and returns messages with WebSocket.
// The returned message will be `<received message> + " ok"`.
type webSocketHandler struct {
	bodies     []string
	times      []time.Time
	milliTimes []int64
}

func (h *webSocketHandler) ServeWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	for {
		msg := ""
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		h.bodies = append(h.bodies, msg)
		now := time.Now()
		h.times = append(h.times, now)
		h.milliTimes = append(h.milliTimes, now.UnixMilli())

		err = websocket.Message.Send(ws, msg+" ok")
		if err != nil {
			panic(err)
		}
	}

}

// receiveBinaryHandler receives binary request body and calculate its md5 hash.
type receiveBinaryHandler struct {
	md5Hash []byte
	size    int64
}

func (h *receiveBinaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := md5.New()
	h.size, _ = io.Copy(hash, r.Body)
	h.md5Hash = hash.Sum(nil)
	w.Write([]byte("ok"))
}

// returnBinaryHandler returns binary body and calculate its md5 hash.
type returnBinaryHandler struct {
	md5Hash []byte
	size    int64
}

func (h *returnBinaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := md5.New()
	hashWriter := io.MultiWriter(hash, w)
	buf := make([]byte, 1000)
	for i := 0; i < 10_000; i++ {
		rand.Read(buf)
		hashWriter.Write(buf)
	}
	h.size = 1000 * 10_000
	h.md5Hash = hash.Sum(nil)
}
