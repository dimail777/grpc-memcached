package memcache

import (
	"fmt"
	"time"
)

type Client interface {
	Get(key string, deadline time.Time) (string, error)
	Set(key, value string, deadline time.Time) error
	Del(key string, deadline time.Time) error
}

type ConnConfig struct {
	Host       string
	Port       int32
	TcpTimeout time.Duration
}

type PoolConfig struct {
	MaxConnections    int
	MaxWaitConnection time.Duration
}

type internalClient struct {
	config ConnConfig
	pool   ConnectionPool
}

func InitClient(connConfig ConnConfig, poolConfig PoolConfig) Client {
	pool := initConnectionPool(poolConfig.MaxConnections, poolConfig.MaxWaitConnection)
	return &internalClient{config: connConfig, pool: pool}
}

func (m *internalClient) Get(key string, deadline time.Time) (string, error) {
	err := isKeyValid(key)
	if err != nil {
		return "", fmt.Errorf("the key is not valid, %s", err)
	}
	before := deadline.Before(time.Now())
	if before {
		return "", fmt.Errorf("deadline has been expired")
	}

	conn, err := m.pool.getConnection(m.config.Host, m.config.Port, m.config.TcpTimeout)
	if err != nil {
		return "", err
	}
	result, err := conn.get(key, deadline)
	if err == nil {
		m.pool.releaseConnection(conn)
	} else {
		m.pool.releaseFailedConnection(conn)
	}
	return result, err
}

func (m *internalClient) Set(key, value string, deadline time.Time) error {
	err := isKeyValid(key)
	if err != nil {
		return fmt.Errorf("the key is not valid, %s", err)
	}
	before := deadline.Before(time.Now())
	if before {
		return fmt.Errorf("deadline has been expired")
	}

	conn, err := m.pool.getConnection(m.config.Host, m.config.Port, m.config.TcpTimeout)
	if err != nil {
		return err
	}
	err = conn.set(key, value, deadline)
	if err == nil {
		m.pool.releaseConnection(conn)
	} else {
		m.pool.releaseFailedConnection(conn)
	}
	return err
}

func (m *internalClient) Del(key string, deadline time.Time) error {
	err := isKeyValid(key)
	if err != nil {
		return fmt.Errorf("the key is not valid, %s", err)
	}
	before := deadline.Before(time.Now())
	if before {
		return fmt.Errorf("deadline has been expired")
	}

	conn, err := m.pool.getConnection(m.config.Host, m.config.Port, m.config.TcpTimeout)
	if err != nil {
		return err
	}
	err = conn.del(key, deadline)
	if err == nil {
		m.pool.releaseConnection(conn)
	} else {
		m.pool.releaseFailedConnection(conn)
	}
	return err
}
