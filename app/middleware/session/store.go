// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package session

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/aileron-gateway/aileron-gateway/util/session"
)

// noopFunc is the function that do nothing.
// This function is used as finish function when tracer is nil.
var noopFunc = func() {}

// cookieSessionStore uses HTTP Cookies as session store.
// Data is split into multiple parts to prevent exceeding size limit of headers.
// Cookie store does not use session ID because the session data itself saved in the cookie after hashed and encrypted hopefully.
// Be careful not to save the cookies in the session because it double and double the session data and reach the header size limit.
// This implements sessionStore interface.
type cookieSessionStore struct {
	// cc is the CookieCreator used to save the session data.
	// Cookie attributes should be set from the stand point of security.
	cc core.CookieCreator

	// enc is the secure encoder.
	// This is used to hashing and encrypting the session data.
	// Both hashing and encrypting can be disabled but is not recommended.
	enc *security.SecureEncoder

	cookiePrefix string

	sm session.SerializeMethod
}

// initialize initializes the cookie store.
// We do not need to care about session ID because the cookie store does not use session ID.
// Overview of initializing flow is as follows.
//
//	1 Create a new session object with a new session ID.
//	2 Get the session data from cookie. Use newly created session if there does not exist.
//	3 Decode the session data. Use newly created session if decode failed.
//	4 Unmarshal the decoded session data. Return an error if failed.
func (h *cookieSessionStore) Get(r *http.Request) (session.Session, error) {
	ss := session.NewDefaultSession(h.sm)

	// Get the session data from the cookie.
	cks := r.Cookies()
	rawData := utilhttp.GetCookie(h.cookiePrefix, r.Cookies())
	if rawData == "" {
		return ss, nil
	}

	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		// The data saved in the cookie might be changed by the user or broken by some network error.
		// That should not be a server side error but client side error.
		// So, use newly created session instead.
		return ss, err
	}

	decoded, err := h.enc.Decode(data)
	if err != nil {
		// The data saved in the cookie might be changed by the user or broken by some network error.
		// That should not be a server side error but client side error.
		// So, use newly created session instead.
		return ss, err
	}

	if err := ss.UnmarshalBinary(decoded); err != nil {
		return nil, err
	}

	ss.Attributes()["__cookie_names"] = utilhttp.CookieNames(cks)
	return ss, nil
}

// save saves the session data into the cookie.
// We do not need to care about session ID because the cookie store does not use session ID.
// Overview of saving flow is as follows.
//
//	1 Delete existing session data if required. Return immediately after deleted.
//	2 Return immediately when session data was not changed.
//	3 Encode the session data. Return an error if failed.
//	4 Save encoded session data into the cookie.
func (h *cookieSessionStore) Save(_ context.Context, w http.ResponseWriter, ss session.Session) error {
	names, _ := ss.Attributes()["__cookie_names"].([]string)

	if ss.SetFlag(0)&session.Delete > 0 { // Should delete session.
		utilhttp.DeleteCookie(w, names, h.cookiePrefix)
		return nil
	}

	data, err := ss.MarshalBinary()
	if err != nil {
		return err
	}

	encoded, err := h.enc.Encode(data)
	if err != nil {
		return err
	}

	value := base64.StdEncoding.EncodeToString(encoded)
	utilhttp.SetCookie(w, names, h.cc, h.cookiePrefix, value)
	return nil
}

type sessionKVS interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte) error
	Delete(context.Context, string) error
}

// kvsSessionStore uses any key-value store as session store.
// Session ID is saved in the cookie and the session data is saved in the key-value store.
// This implements sessionStore interface.
type kvsSessionStore struct {
	store      sessionKVS
	enc        *security.SecureEncoder
	cc         core.CookieCreator
	tracer     app.Tracer
	cookieName string
	prefix     string
	sm         session.SerializeMethod
}

// initialize initializes the session store.
// Overview of initializing flow is as follows.
//
//	1 Create a new session object with a new session ID.
//	2 Get session ID from cookie. Use newly created session if there does not exist.
//	3 Get the session data from the key-value store. Use newly created session if there does not exist.
//	4 Decode the session data. Use newly created session if decode failed.
//	5 Unmarshal the decoded session data. Return an error if failed.
//	6 Replace the session ID with the existing one.
func (s *kvsSessionStore) Get(r *http.Request) (session.Session, error) {
	ss := session.NewDefaultSession(s.sm)

	ck, err := r.Cookie(s.cookieName) // Session ID from cookie.
	if err == http.ErrNoCookie {
		return ss, err
	}
	ssid := ck.Value

	// if !session.Validator.MatchString(ssid) {
	// 	// Invalid format of session ID detected.
	// 	// Session ID might be intentionally or accidentally changed by clients.
	// 	// So, newly created session with its new session ID.
	// 	return ss, nil
	// }

	ctx := r.Context()
	finish := noopFunc
	if s.tracer != nil {
		ctx, finish = s.tracer.Trace(ctx, "session", map[string]string{"operation": "GET"})
	}
	data, err := s.store.Get(ctx, s.prefix+ssid)
	finish() // Finish tracing span.

	if err != nil {
		if err == kvs.Nil {
			// Session ID exists in the cookie.
			// But session data was not found in the store.
			return ss, nil
		} else {
			// Some error occurred while operating the key-value storage.
			// Do not delete the session data because
			// the error may caused by a networking.
			return nil, err
		}
	}

	decoded, err := s.enc.Decode(data)
	if err != nil {
		return ss, nil
	}

	if err := ss.UnmarshalBinary(decoded); err != nil {
		return nil, err
	}

	ss.Attributes()["__session_id"] = ssid
	return ss, nil
}

// save saves the session data into store.
// Overview of saving flow is as follows.
//
//	1 Delete existing session data if required. Return immediately after deleted.
//	2 Return immediately when session data was not changed or there is no need to refresh session ID.
//	3 Encode the session data. Return an error if failed.
//	4 Save encoded session data into the cookie.
//	5 Set session ID in the cookie when not set or session ID refreshed.
func (s *kvsSessionStore) Save(ctx context.Context, w http.ResponseWriter, ss session.Session) error {
	ssid, _ := ss.Attributes()["__session_id"].(string)
	flags := ss.SetFlag(0)

	if flags&session.Delete > 0 {
		if flags&session.Restored > 0 {
			return nil
		}
		finish := noopFunc
		if s.tracer != nil {
			ctx, finish = s.tracer.Trace(ctx, "session", map[string]string{"operation": "DELETE"})
		}
		s.store.Delete(ctx, s.prefix+ssid)
		finish()
		http.SetCookie(w, &http.Cookie{Name: s.cookieName, MaxAge: -1})
		return nil
	}

	if flags&session.Refresh > 0 {
		b, err := uid.NewHostedID()
		if err != nil {
			return err
		}
		ssid = base32.HexEncoding.EncodeToString(b)
	}

	// Get a byte slice of the session.
	data, err := ss.MarshalBinary()
	if err != nil {
		return err
	}

	encoded, err := s.enc.Encode(data)
	if err != nil {
		return err
	}

	// Use tracer if it is given.
	finish := noopFunc
	if s.tracer != nil {
		ctx, finish = s.tracer.Trace(ctx, "session", map[string]string{"operation": "SET"})
	}
	err = s.store.Set(ctx, s.prefix+ssid, encoded)
	finish() // Finish tracing span.
	if err != nil {
		return err
	}

	ck := s.cc.NewCookie()
	ck.Name = s.cookieName
	ck.Value = ssid // Set new session ID.
	http.SetCookie(w, ck)

	return nil
}
