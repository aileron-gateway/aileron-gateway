package session

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

// NoValue means no value in the session.
//
//nolint:staticcheck // ST1012: error var Nil should have name of the form ErrFoo
//lint:ignore ST1012 // error var Nil should have name of the form ErrFoo
var NoValue = errors.New("util/session: no value")

// SerializeMethod is a session data serialization method.
// Currently available:
//   - SerializeMsgPack
//   - SerializeJSON
type SerializeMethod int

const (
	SerializeMsgPack SerializeMethod = iota
	SerializeJSON
)

const (
	Noop     uint = 1 << iota // No meaning.
	New                       // Session is New.
	Restored                  // Session is restored from session store.
	Updated                   // Session data is updated.
	Delete                    // Delete the session.
	Refresh                   // Refresh the session ID.
)

// Store is the session store.
type Store interface {
	// Get returns a new session.
	// Nil session and non-nil error should be server-side error.
	// Non-nil session an non-nil error should be client-side error.
	// If client-error occurred, callers can use the returned session
	// as a new session. Implementers should return a new empty session.
	Get(r *http.Request) (Session, error)

	// Save saves session data.
	// Non-nil error should be server-side error.
	Save(context.Context, http.ResponseWriter, Session) error
}

// Session is the interface for session objects.
type Session interface {
	// BinaryMarshaler marshals session data.
	// Attributes are not marshaled.
	encoding.BinaryMarshaler
	// BinaryUnmarshaler unmarshals given data
	// into this session object.
	encoding.BinaryUnmarshaler

	// SetFlag sets a flag to this object
	// and returns the resulting flag sets.
	// Once flag  set, it can not be unset.
	// flag=0 has takes no effects.
	SetFlag(flag uint) (flags uint)

	// Attributes is the additional data bounded
	// to this session object.
	// Attributes are not persistent.
	Attributes() map[string]any

	// Delete deletes the data with given key
	// from this session object.
	Delete(key string)

	// Persist persists the given object.
	// If the given value implements encoding.BinaryMarshaler
	// interface, its method is called for marshaling.
	Persist(key string, value any) error

	// Extract extracts data into given object.
	// If the given object implements encoding.BinaryUnmarshaler
	// interface, its method is called for unmarshaling.
	Extract(key string, into any) error
}

func NewDefaultSession(sm SerializeMethod) *DefaultSession {
	df := &DefaultSession{
		flags: New,
		attrs: map[string]any{},
		data:  map[string][]byte{},
	}
	if sm == SerializeJSON {
		df.marshal, df.unmarshal = json.Marshal, json.Unmarshal
	} else {
		df.marshal, df.unmarshal = msgpack.Marshal, msgpack.Unmarshal
	}
	return df
}

type DefaultSession struct {
	raw       []byte
	data      map[string][]byte
	attrs     map[string]any
	flags     uint
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (s *DefaultSession) SetFlag(flag uint) uint {
	s.flags |= flag
	return s.flags
}

func (s *DefaultSession) Attributes() map[string]any {
	return s.attrs
}

func (s *DefaultSession) Delete(key string) {
	if _, ok := s.data[key]; ok {
		s.flags |= Updated
		delete(s.data, key)
	}
}

func (s *DefaultSession) Persist(key string, value any) error {
	var b []byte
	var err error
	if bm, ok := value.(encoding.BinaryMarshaler); ok {
		b, err = bm.MarshalBinary()
	} else {
		b, err = s.marshal(value)
	}
	if err != nil {
		return err
	}
	s.data[key] = b
	s.flags |= Updated
	return nil
}

func (s *DefaultSession) Extract(key string, into any) error {
	b, ok := s.data[key]
	if !ok {
		return NoValue
	}
	if bu, ok := into.(encoding.BinaryUnmarshaler); ok {
		return bu.UnmarshalBinary(b)
	} else {
		return s.unmarshal(b, into)
	}
}

func (s *DefaultSession) UnmarshalBinary(raw []byte) error {
	s.flags |= Restored
	s.flags &= (^New)
	s.raw = raw
	return s.unmarshal(raw, &s.data)
}

func (s *DefaultSession) MarshalBinary() ([]byte, error) {
	if s.flags&Updated == 0 {
		return s.raw, nil
	}
	return s.marshal(&s.data)
}

func MustPersist(ss Session, key string, value any) {
	if err := ss.Persist(key, value); err != nil {
		panic(err)
	}
}
