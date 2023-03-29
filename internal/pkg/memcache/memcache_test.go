package memcache

import (
	"context"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithMemCached(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "memcached",
		Hostname:     "127.0.0.1",
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	memCacheC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start memcahced, %s", err)
	}
	portNat, err := memCacheC.MappedPort(ctx, "11211/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %s", err.Error())
	}
	port, err := convertPortToInt(portNat)
	if err != nil {
		t.Fatalf("failed to convert port: %s", err.Error())
	}

	defer func() {
		if err := memCacheC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	testMemCacheConnection(t, port)
	testPoolConnection(t, port)
	testClient(t, port)
}

func convertPortToInt(portNat nat.Port) (int32, error) {
	onlyPort := strings.Split(string(portNat), "/")[0]
	port, err := strconv.ParseUint(onlyPort, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(port), nil
}

func testMemCacheConnection(t *testing.T, port int32) {
	conn, err := initConnection("127.0.0.1", port, 5*time.Second)
	if err != nil {
		t.Fatalf("connection error, %s", err)
	}
	value, err := conn.get("tes_key", time.Now().Add(5*time.Second))
	if err != nil {
		t.Fatalf("get command error, %s", err)
	}
	if value != "" {
		t.Fatalf("get value error, %s != empty", value)
	}

	err = conn.set("tes_key", "test_value", time.Now().Add(5*time.Second))
	if err != nil {
		t.Fatalf("set command error, %s", err)
	}

	value, err = conn.get("tes_key", time.Now().Add(5*time.Second))
	if err != nil {
		t.Fatalf("get command error, %s", err)
	}
	if value != "test_value" {
		t.Fatalf("get value error, %s != test_value", value)
	}

	err = conn.del("tes_key", time.Now().Add(5*time.Second))
	if err != nil {
		t.Fatalf("del command error, %s", err)
	}

	value, err = conn.get("tes_key", time.Now().Add(5*time.Second))
	if err != nil {
		t.Fatalf("get command error, %s", err)
	}
	if value != "" {
		t.Fatalf("get value error, %s != empty", value)
	}
}

func testPoolConnection(t *testing.T, port int32) {
	pool := initConnectionPool(2, 1*time.Second)
	conn, err := pool.getConnection("127.0.0.1", port, 1*time.Second)
	if err != nil {
		t.Fatalf("connection error, %s", err)
	}
	pool.releaseConnection(conn)
	conn, err = pool.getConnection("127.0.0.1", port, 1*time.Second)
	if err != nil {
		t.Fatalf("connection error, %s", err)
	}
	pool.releaseFailedConnection(conn)

	conn1, _ := pool.getConnection("127.0.0.1", port, 1*time.Second)
	conn2, _ := pool.getConnection("127.0.0.1", port, 1*time.Second)
	_, err = pool.getConnection("127.0.0.1", port, 1*time.Second)
	if err == nil {
		t.Fatalf("connection wait error, 1s")
	}
	pool.releaseConnection(conn1)
	pool.releaseConnection(conn2)
	conn1, err = pool.getConnection("127.0.0.1", port, 1*time.Second)
	if err != nil {
		t.Fatalf("connection error, %s", err)
	}
	pool.releaseConnection(conn1)
}

func testClient(t *testing.T, port int32) {
	connConfig := ConnConfig{
		Host:       "127.0.0.1",
		Port:       port,
		TcpTimeout: 1 * time.Second,
	}
	poolConfig := PoolConfig{
		MaxConnections:    2,
		MaxWaitConnection: 1 * time.Second,
	}
	client := InitClient(connConfig, poolConfig)

	parallel := 50
	callback := make(chan error, parallel)
	for i := 0; i < parallel; i++ {
		go testClientInParallel(client, callback)
	}
	for i := 0; i < parallel; i++ {
		res := <-callback
		if res != nil {
			t.Fatalf("test e2e in parallel error, %s", res)
		}
	}
}

func testClientInParallel(client Client, callback chan error) {
	err := client.Set("key", "value_value_value_value_value_value_value", time.Now().Add(1*time.Second))
	if err != nil {
		callback <- err
		return
	}
	_, err = client.Get("key", time.Now().Add(1*time.Second))
	if err != nil {
		callback <- err
		return
	}
	err = client.Del("key", time.Now().Add(1*time.Second))
	if err != nil {
		callback <- err
		return
	}
	callback <- nil
}
