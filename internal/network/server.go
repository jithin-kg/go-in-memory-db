package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/jithin-kg/go-in-memory-db/internal/db"
	"github.com/jithin-kg/go-in-memory-db/internal/protocol"
)

// Conn is an interface that abstracts the methods from net.Conn used in handleConnection
type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// Use Go’s net package to listen for connections and spawn a new
//  goroutine for each connection to handle requests

func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %s, err: %v", port, err)
	}
	defer listener.Close()
	log.Printf("Server listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection %v", err)
		}
		go HandleConnection(conn)
	}
}

var kvStore = db.NewKeyValueStore()

func sendResponse(conn net.Conn, status byte, data []byte) {
	// prepare the respose header
	responseHeader := []byte{status}

	if len(data) > 0 {
		// include the data length there is some data to be send
		dataLength := make([]byte, 2)
		binary.BigEndian.PutUint16(dataLength, uint16(len(data)))
		responseHeader = append(responseHeader, dataLength...)
		responseHeader = append(responseHeader, data...)
	} else {
		responseHeader = append(responseHeader, 0x00, 0x00)
	}
	// send the response
	if _, err := conn.Write(responseHeader); err != nil {
		log.Printf("Failed to send response: %v\n", err)
	}
}

// func handleConnection(conn net.Conn) {
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		opCode, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				return
			}
			// data -> *errors.errorString {s: "EOF"}
			sendResponse(conn, 0xFF, []byte(fmt.Sprintf("Error reading operation code %v", err)))
			log.Println("Error reading opCode:", err)
			return
		}
		keyLengthBuf := make([]byte, 2)
		_, err = reader.Read(keyLengthBuf)
		if err != nil {
			sendResponse(conn, 0xFF, []byte(fmt.Sprintf("Error reading key length: %v", err)))
			log.Println("Error reading key length:", err)
			return
		}
		keyLength := binary.BigEndian.Uint16(keyLengthBuf) // converting []byte of length 2 into unsigned integer
		// eg: converts [0x00, 0x03] -> 3
		keyBuf := make([]byte, keyLength)
		_, err = reader.Read(keyBuf)

		if err != nil {
			sendResponse(conn, 0xFF, []byte(fmt.Sprintf("Error reading key: %v", err)))
			log.Println("Error reading key:", err)
			return
		}
		switch opCode {
		// Only process value for SET operation (0x01)
		case protocol.OpSet:
			valueLengthBuf := make([]byte, 2)
			_, err := reader.Read(valueLengthBuf)
			if err != nil {
				sendResponse(conn, 0xFF, []byte(fmt.Sprintf("Error reading value length: %v", err)))
				log.Println("Error reading value length:", err)
				return
			}
			valueLength := binary.BigEndian.Uint16(valueLengthBuf)

			valueBuf := make([]byte, valueLength)
			_, err = reader.Read(valueBuf)
			if err != nil {
				sendResponse(conn, 0xFF, []byte(fmt.Sprintf("Error reading value: %v", err)))
				log.Println("Error reading value:", err)
				return
			}
			key := string(keyBuf)
			value := string(valueBuf)
			kvStore.Set(key, value)
			sendResponse(conn, 0x00, []byte("SET successful"))
			log.Printf("SET operation for key: %s, value: %s\n", keyBuf, valueBuf)
		case protocol.OpGet:
			// handle the GET operation
			log.Printf("GET operation for key: %s\n", keyBuf)
		}

	}
}
