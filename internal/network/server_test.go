package network_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jithin-kg/go-in-memory-db/internal/network"
	"github.com/jithin-kg/go-in-memory-db/internal/protocol"
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

func newMockConnection(input, output *bytes.Buffer, localAddr, remoteAddr net.Addr) network.Conn {
	return &mockConn{
		readBuffer:  input,
		writeBuffer: output,
		localAddr:   localAddr,
		remoteAddr:  remoteAddr,
	}
}
func TestHandleConnection(t *testing.T) {
	// simulate SET operation request
	input := bytes.NewBuffer([]byte{
		protocol.OpSet, // op code for SET
		0x00, 0x03,     //key length (3)
		'k', 'e', 'y',
		0x00, 0x05, // value length (5)
		'v', 'a', 'l', 'u', 'e',
	})
	output := new(bytes.Buffer)
	localAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	remoteAddr := &net.TCPAddr{IP: net.ParseIP("10.0.0.1"), Port: 9090}
	conn := newMockConnection(input, output, localAddr, remoteAddr)
	network.HandleConnection(conn)

	// Read and validate the response
	expectedStatus := byte(0x00) // 0x00 indicates success
	status := output.Bytes()[0]

	if status != expectedStatus {
		t.Errorf("Expected status %v, got: %v", expectedStatus, status)
	}
	// next two bytes indicates the length of the message
	// converting it into uint16
	dataLength := binary.BigEndian.Uint16(output.Bytes()[1:3])
	if dataLength > 0 {
		message := string(output.Bytes()[3 : 3+dataLength])
		fmt.Printf("Response message %s", message)
	} else {
		t.Errorf("Expected non-zero data length in resonse")
	}
}
