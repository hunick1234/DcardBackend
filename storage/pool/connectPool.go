package pool

import (
	"log"
	"sync"
	"time"

	"github.com/hunick1234/DcardBackend/config"
	"github.com/hunick1234/DcardBackend/storage"
)

type Pool struct {
	dbConnections sync.Map
}

func NewPool() *Pool {
	return &Pool{
		dbConnections: sync.Map{},
	}
}

func (p *Pool) setDBConnection(dbName string, conn storage.Storager) {
	p.dbConnections.Store(dbName, conn)
}

func (p *Pool) getDBConnection(dbName string) (storage.Storager, bool) {
	storager, ok := p.dbConnections.Load(dbName)
	if !ok {
		return nil, false
	}
	return storager.(storage.Storager), true
}

func (p *Pool) deleteDBConnection(dbName string) {
	p.dbConnections.Delete(dbName)
}

func (p *Pool) rangeDBConnection(f func(dbName string, conn storage.Storager) bool) {
	p.dbConnections.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(storage.Storager))
	})
}

func (p *Pool) GetConnection(cfg *config.MongoCfg) (storage.Storager, error) {
	start := time.Now()
	defer func() {
		log.Println("get connection time: ", time.Since(start))
	}()

	var err error

	if conn, ok := p.getDBConnection(cfg.DB); ok {
		if ChcekConn(conn) {
			return conn, nil
		}
		p.deleteDBConnection(cfg.DB)
	}

	//if not exist, create new connection

	//if not exist, create new connection
	storeger, err := storage.NewMongoConn(cfg)
	if err != nil {
		return nil, err
	}

	p.setDBConnection(cfg.DB, storeger)

	return storeger, nil
}

func (p *Pool) ClosePool() {
	p.rangeDBConnection(func(dbName string, conn storage.Storager) bool {
		conn.Disconnect()
		return true
	})
}

func (p *Pool) Disconnect(conn storage.Storager) {
	err := conn.Disconnect()
	if err != nil {
		return
	}

	p.dbConnections.Delete(conn.GetDBName())
}
func ChcekConn(conn storage.Storager) bool {
	if conn.Ping() != nil {
		return false
	}
	return true
}
