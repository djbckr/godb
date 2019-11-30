package session

import "time"

type Session struct {
	lastActive  time.Time     // the last time this object was doing something
	initialized time.Time     // when this session was started
	username    string        // the username for this session
	authKey     int           // the authorization key for this session - there can be more than one session with the same authKey
	sessionId   int           // the unique identifier for this session
	maxIdleTime time.Duration // the allowed idle time for this session
}

var sessionListId map[int]*Session

func init() {
	sessionListId = make(map[int]*Session)
}

func SessionById(id int) *Session {
	session := sessionListId[id]

	if session == nil {
		return nil
	}

	session.lastActive = time.Now()

	return session
}
