// Package srv implements the server components required to manage the database
// and its keys.
package srv

import "sync"

// Server defines the methods required to implement a Server.
type Server interface {
	RWLocker

	GetDB(int) (DB, bool)
	FlushAll()
}

const maxDBs = 16 // TODO : Should be read from configuration

// Static check to make sure *server implements the Server interface.
var _ Server = (*server)(nil)

// The one and only server instance.
var DefaultServer Server

// server is the internal implementation of a Server.
type server struct {
	sync.RWMutex
	dbs []DB
}

func init() {
	// TODO : Read configuration
	DefaultServer = &server{
		dbs: make([]DB, maxDBs),
	}
}

// FlushAll clears the keys from all databases.
func (s *server) FlushAll() {
	s.dbs = make([]DB, maxDBs)
}

// GetDB returns the database identified by its index.
func (s *server) GetDB(ix int) (DB, bool) {
	if ix < 0 || ix >= maxDBs {
		return nil, false
	}

	db := s.dbs[ix]
	if db == nil {
		db = NewDB(ix)
		s.dbs[ix] = db
	}
	return db, true
}
