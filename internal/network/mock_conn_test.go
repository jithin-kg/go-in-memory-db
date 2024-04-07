package network_test

import (
	"bytes"
	"net"
	"time"
)

type mockConn struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	localAddr   net.Addr
	remoteAddr  net.Addr
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}
func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
