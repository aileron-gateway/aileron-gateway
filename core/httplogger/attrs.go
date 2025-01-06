package httplogger

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	keyType     = "type"
	keyID       = "id"
	keyTime     = "time"
	keyMethod   = "method"
	keyRemote   = "remote"
	keyHost     = "host"
	keyPath     = "path"
	keyQuery    = "query"
	keySize     = "size"
	keyStatus   = "status"
	keyDuration = "duration"
	keyProto    = "proto"
	keyHeader   = "header"
	keyBody     = "body"
	keyRead     = "read"
	keyWritten  = "written"
)

type requestAttrs struct {
	typ    string            // Typ is the log type "server" or "client".
	id     string            // id is the log ID in the case for there is no request ID.
	time   string            // time is the observed request time.
	host   string            // host is the host name which clients requested with.
	method string            // method is the requested HTTP method.
	path   string            // path is the requested URL path.
	query  string            // query is the requested query URa.
	remote string            // remote is the remote, or client address.
	proto  string            // proto is the protocol version.
	size   int64             // body size in bytes.
	header map[string]string // headers is the all request headers.
	body   string            // body is the request body.
}

func (a *requestAttrs) accessKeyValues() []any {
	return []any{
		"request",
		map[string]any{
			keyID:     a.id,
			keyTime:   a.time,
			keyHost:   a.host,
			keyMethod: a.method,
			keyPath:   a.path,
			keyQuery:  a.query,
			keyRemote: a.remote,
			keyProto:  a.proto,
			keySize:   a.size,
			keyHeader: a.header,
		},
	}
}

func (a *requestAttrs) journalKeyValues() []any {
	return []any{
		"request",
		map[string]any{
			keyID:     a.id,
			keyTime:   a.time,
			keyHost:   a.host,
			keyMethod: a.method,
			keyPath:   a.path,
			keyQuery:  a.query,
			keyRemote: a.remote,
			keyProto:  a.proto,
			keySize:   a.size,
			keyHeader: a.header,
			keyBody:   a.body,
		},
	}
}

func (a *requestAttrs) TagFunc(tag string) []byte {
	switch tag {
	case keyID:
		return []byte(a.id)
	case keyTime:
		return []byte(a.time)
	case keyMethod:
		return []byte(a.method)
	case keyPath:
		return []byte(a.path)
	case keyQuery:
		return []byte(a.query)
	case keyHost:
		return []byte(a.host)
	case keyRemote:
		return []byte(a.remote)
	case keyProto:
		return []byte(a.proto)
	case keySize:
		return []byte(strconv.FormatInt(a.size, 10))
	case keyHeader:
		return []byte(fmt.Sprint(a.header)) // "%+v"
	case keyBody:
		return []byte(a.body)
	case keyType:
		return []byte(a.typ)
	default:
		switch {
		case strings.HasPrefix(tag, "header."):
			return []byte(a.header[textproto.CanonicalMIMEHeaderKey(tag[7:])])
		}
		return []byte("<undefined:" + tag + ">")
	}
}

type responseAttrs struct {
	typ      string            // Typ is the log type "server" or "client".
	id       string            // id is the log ID in the case for there is no request ID.
	time     string            // time is the observed request time.
	duration int64             // duration of the request.
	status   int               // status is the response status.
	size     int64             // body size.
	header   map[string]string // headers is all response headers.
	body     string            // body is response body.
}

func (a *responseAttrs) accessKeyValues() []any {
	return []any{
		"response",
		map[string]any{
			keyID:       a.id,
			keyTime:     a.time,
			keyDuration: a.duration,
			keyStatus:   a.status,
			keySize:     a.size,
			keyHeader:   a.header,
		},
	}
}

func (a *responseAttrs) journalKeyValues() []any {
	return []any{
		"response",
		map[string]any{
			keyID:       a.id,
			keyTime:     a.time,
			keyDuration: a.duration,
			keyStatus:   a.status,
			keySize:     a.size,
			keyHeader:   a.header,
			keyBody:     a.body,
		},
	}
}

func (a *responseAttrs) TagFunc(tag string) []byte {
	switch tag {
	case keyID:
		return []byte(a.id)
	case keyTime:
		return []byte(a.time)
	case keyDuration:
		return []byte(strconv.FormatInt(a.duration, 10))
	case keyStatus:
		return []byte(strconv.FormatInt(int64(a.status), 10))
	case keySize:
		return []byte(strconv.FormatInt(a.size, 10))
	case keyHeader:
		return []byte(fmt.Sprint(a.header)) // "%+v"
	case keyBody:
		return []byte(a.body)
	case keyType:
		return []byte(a.typ)
	default:
		switch {
		case strings.HasPrefix(tag, "header."):
			return []byte(a.header[textproto.CanonicalMIMEHeaderKey(tag[7:])])
		}
		return []byte("<undefined:" + tag + ">")
	}
}

// idContext is the context key type to save
// log IDs in contexts.
type idContext struct{}

var (
	// idContextKey is the context key to save
	// log IDs in contexts.
	idContextKey = &idContext{}

	// counter is the atomic and cyclic counter that is
	// added on every requests.
	counter atomic.Uint64

	// hostname is the FNV1_32a hashed hostname.
	// This value is used for generating log IDs.
	hostname = func() []byte {
		hostname, _ := os.Hostname()
		h := fnv.New32a()
		h.Write([]byte(hostname))
		return h.Sum(nil)
	}()
)

// newLogID returns a new log ID string.
// It consists of
//   - 4 bytes (32 bits) hostname FNV1a hash
//   - 4 bytes (32 bits) unix seconds which cycle every 4,294,967,295 seconds or 136 years.
//   - 7 bytes (56 bits) counter which cycles every 7,2057,594,037,927,935 or 0xFFFFFFFFFFFFFF.
func newLogID() string {
	x := [15]byte{}
	copy(x[:], hostname)

	c := counter.Add(1)
	x[14] = byte(c)
	x[13] = byte(c >> 8)
	x[12] = byte(c >> 16)
	x[11] = byte(c >> 24)
	x[10] = byte(c >> 32)
	x[9] = byte(c >> 40)
	x[8] = byte(c >> 48)

	s := time.Now().Unix()
	x[7] = byte(s)
	x[6] = byte(s >> 8)
	x[5] = byte(s >> 16)
	x[4] = byte(s >> 24)

	return base64.URLEncoding.EncodeToString(x[:])
}
