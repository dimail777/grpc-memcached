package memcache

import (
	"fmt"
	"sync"
	"time"
)

type ConnectionPool interface {
	getConnection(host string, port int32, tcpTimeout time.Duration) (Connection, error)
	releaseConnection(conn Connection)
	releaseFailedConnection(conn Connection)
}

type internalConnectionPool struct {
	mu                 sync.Mutex
	createdConnections int
	maxConnections     int
	freeConnections    chan Connection
	failedConnections  chan Connection
	maxWaitConnection  time.Duration
}

func initConnectionPool(maxConnections int, maxWaitConnection time.Duration) ConnectionPool {
	return &internalConnectionPool{
		createdConnections: 0,
		maxConnections:     maxConnections,
		maxWaitConnection:  maxWaitConnection,
		freeConnections:    make(chan Connection, maxConnections),
		failedConnections:  make(chan Connection, maxConnections),
	}
}

func (cp *internalConnectionPool) getConnection(host string, port int32, tcpTimeout time.Duration) (Connection, error) {
	for {
		cp.mu.Lock()
		if cp.createdConnections >= cp.maxConnections {
			cp.mu.Unlock()
			select {
			case connection := <-cp.freeConnections:
				return connection, nil
			case <-cp.failedConnections:
				continue
			case <-time.After(cp.maxWaitConnection):
				return nil, fmt.Errorf("timeout waiting for connection")
			}
		}
		conn, err := initConnection(host, port, tcpTimeout)
		if err == nil {
			cp.createdConnections++
		}
		cp.mu.Unlock()
		return conn, err
	}
}

func (cp *internalConnectionPool) releaseConnection(conn Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.freeConnections <- conn
}

func (cp *internalConnectionPool) releaseFailedConnection(conn Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	_ = conn.close()
	cp.failedConnections <- conn
	cp.createdConnections--
}
