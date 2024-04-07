package network_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/jithin-kg/go-in-memory-db/internal/network"
	"github.com/jithin-kg/go-in-memory-db/internal/protocol"
)

type MockDbService struct {
	store map[string]interface{}
}

func (s *MockDbService) Set(key string, value interface{}) error {
	s.store[key] = value
	return nil
}

func (s *MockDbService) Get(key string) (interface{}, bool, error) {
	val, ok := s.store[key]
	return val, ok, nil
}
func NewMockDbService() *MockDbService {
	return &MockDbService{
		store: make(map[string]interface{}),
	}
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

	mockService := NewMockDbService()
	server := network.NewServer(mockService)
	server.HandleConnection(conn)
	expectedStatus := byte(0x00) // 0x00 indicates success

	validateResponse(t, output, expectedStatus, network.STATUS_SET_SUCCESS)

}

func TestOperations(t *testing.T) {
	mockService := NewMockDbService()
	server := network.NewServer(mockService)
	localAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	remoteAddr := &net.TCPAddr{IP: net.ParseIP("10.0.0.1"), Port: 9090}

	tests := []struct {
		name  string
		input *bytes.Buffer
	}{
		{
			name: "Test SET operation",
			input: bytes.NewBuffer([]byte{
				protocol.OpSet, // op code for SET
				0x00, 0x03,     //key length (3)
				'k', 'e', 'y',
				0x00, 0x005, //valu length (5)
				'v', 'a', 'l', 'u', 'e',
			}),
		}, {
			name: "Test GET opeartion",
			input: bytes.NewBuffer([]byte{
				protocol.OpGet,
				0x00, 0x003, //key length
				'k', 'e', 'y',
			}),
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			output := new(bytes.Buffer)
			conn := newMockConnection(tst.input, output, localAddr, remoteAddr)
			server.HandleConnection(conn)
			log.Println(output)
			expectedStatus := byte(0x00) // 0x00 indicates success
			validateResponse(t, output, expectedStatus, network.STATUS_SET_SUCCESS)
		})
	}

}

func validateResponse(t *testing.T, output *bytes.Buffer, expectedStatus byte, expectedMessage string) {
	t.Helper()
	// Read and validate the response
	status := output.Bytes()[0]

	if status != expectedStatus {
		t.Errorf("Expected status %v, got: %v", expectedStatus, status)
	}
	// next two bytes indicates the length of the message
	// converting it into uint16
	dataLength := binary.BigEndian.Uint16(output.Bytes()[1:3])
	if dataLength > 0 {
		message := string(output.Bytes()[3 : 3+dataLength])
		if expectedMessage != message {
			t.Errorf("Expected message %v, got: %v", expectedMessage, message)
		}
		fmt.Printf("Response message %s", message)
	} else {
		t.Errorf("Expected non-zero data length in resonse")
	}
}
