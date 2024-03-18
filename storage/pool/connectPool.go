package pool

import (
	"fmt"
	"sync"

	"github.com/hunick1234/DcardBackend/config"
	"github.com/hunick1234/DcardBackend/storage"
)

var mu sync.Mutex

type Pool struct {
	dbConnections map[string]storage.Storager
}

func NewPool() *Pool {
	return &Pool{
		dbConnections: make(map[string]storage.Storager, 10),
	}
}

func (p *Pool) GetConnection(cfg *config.MongoCfg) (storage.Storager, error) {
	var err error
	mu.Lock()
	defer mu.Unlock()
	//if exist return connection
	if p.dbConnections[cfg.DB] != nil {
		if ChcekConn(p.dbConnections[cfg.DB]) {
			return p.dbConnections[cfg.DB], nil
		}
		p.dbConnections[cfg.DB] = nil
	}

	//if not exist, create new connection
	storeger, err := storage.NewMongoConn(cfg)
	if err != nil {
		return nil, err
	}

	p.dbConnections[cfg.DB] = storeger
	return storeger, nil
}

func (p *Pool) ClosePool() {
	for _, conn := range p.dbConnections {
		conn.Disconnect()
	}
}

func (p *Pool) Disconnect(conn storage.Storager) {
	err := conn.Disconnect()
	if err != nil {
		return
	}
	p.dbConnections[conn.GetDBName()] = nil
}
func ChcekConn(conn storage.Storager) bool {
	if conn.Ping() != nil {
		return false
	}
	return true
}

func (p *Pool) LockConnection(dbName string) (storage.Storager, error) {
	mu.Lock()
	defer mu.Unlock()

	conn, exists := p.dbConnections[dbName]
	if !exists || !ChcekConn(conn) {
		return nil, fmt.Errorf("connection does not exist or is not healthy")
	}

	return conn, nil
}

func (p *Pool) UnlockConnection(dbName string) {
	mu.Unlock()
}
