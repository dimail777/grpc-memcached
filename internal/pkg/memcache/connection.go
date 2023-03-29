package memcache

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Connection interface {
	get(key string, deadline time.Time) (string, error)
	set(key, value string, deadline time.Time) error
	del(key string, deadline time.Time) error
	close() error
}

func initConnection(host string, port int32, tcpTimeout time.Duration) (Connection, error) {
	address := host + ":" + strconv.FormatInt(int64(port), 10)
	conn, err := net.DialTimeout("tcp", address, tcpTimeout)
	if err != nil {
		return nil, fmt.Errorf("memcache connection unreachable, %s", err)
	}
	return &internalConnection{
		conn,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	}, nil
}

type internalConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (m *internalConnection) get(key string, deadline time.Time) (string, error) {
	defer m.reader.Reset(m.conn)
	defer m.writer.Reset(m.conn)

	err := m.conn.SetDeadline(deadline)
	if err != nil {
		return "", fmt.Errorf("deadline error, %s", err)
	}

	_, err = m.writer.WriteString("get " + key + "\r\n")
	if err != nil {
		return "", fmt.Errorf("unable to proceed the get command, %s", err)
	}
	err = m.writer.Flush()
	if err != nil {
		return "", fmt.Errorf("unable to proceed the get command, %s", err)
	}
	line, err := m.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("unable to fetch command result, %s", err)
	}
	if isEnd(line) {
		return "", nil
	}
	if !isAValue(line) {
		return "", fmt.Errorf("unexpected error, %s", err)
	}
	size, err := extractValueSize(line)
	if err != nil {
		return "", fmt.Errorf("wrong format of protocol, %s", err)
	}
	size += 7 // with '\r\nEND\r\n' after value
	valueBuf := make([]byte, size)
	readBytes := 0
	for readBytes < size {
		n, err := m.reader.Read(valueBuf[readBytes:])
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("unable to fetch a value, %s", err)
		}
		readBytes += n
	}
	return string(valueBuf[:size-7]), nil
}

func (m *internalConnection) set(key, value string, deadline time.Time) error {
	defer m.reader.Reset(m.conn)
	defer m.writer.Reset(m.conn)
	size := strconv.FormatInt(int64(len(value)), 10)

	err := m.conn.SetDeadline(deadline)
	if err != nil {
		return fmt.Errorf("deadline error, %s", err)
	}

	_, err = m.writer.WriteString("set " + key + " 0 0 " + size + "\r\n")
	if err != nil {
		return fmt.Errorf("unable to proceed the set command, %s", err)
	}
	valueBuf := []byte(value)
	writeBytes := 0
	for writeBytes < len(value) {
		n, err := m.writer.Write(valueBuf[writeBytes:])
		if err != nil && err != io.EOF {
			return fmt.Errorf("unable to load a value, %s", err)
		}
		writeBytes += n
	}
	_, err = m.writer.WriteString("\r\nEND\r\n")
	if err != nil {
		return fmt.Errorf("unable to complete the set command, %s", err)
	}
	err = m.writer.Flush()
	if err != nil {
		return fmt.Errorf("unable to complete the set command, %s", err)
	}

	line, err := m.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("unable to fetch command result, %s", err)
	}
	if isStored(line) {
		return nil
	}
	return fmt.Errorf("unexpected error, %s", err)
}

func (m *internalConnection) del(key string, deadline time.Time) error {
	defer m.reader.Reset(m.conn)
	defer m.writer.Reset(m.conn)

	err := m.conn.SetDeadline(deadline)
	if err != nil {
		return fmt.Errorf("deadline error, %s", err)
	}

	_, err = m.writer.WriteString("delete " + key + "\r\n")
	if err != nil {
		return fmt.Errorf("unable to proceed the del command, %s", err)
	}
	err = m.writer.Flush()
	if err != nil {
		return fmt.Errorf("unable to proceed the del command, %s", err)
	}
	line, err := m.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("unable to fetch command result, %s", err)
	}
	if isNotFount(line) || isDeleted(line) {
		return nil
	}
	return fmt.Errorf("unexpected error, %s", err)
}

func (m *internalConnection) close() error {
	err := m.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
