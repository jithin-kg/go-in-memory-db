package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jithin-kg/go-in-memory-db/internal/db"
	"github.com/jithin-kg/go-in-memory-db/internal/protocol"
)

const (
	STATUS_SET_SUCCESS = "SET successful"
)

type Server struct {
	dbService db.Service
}

func NewServer(dbService db.Service) *Server {
	return &Server{
		dbService: dbService,
	}
}
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

//	goroutine for each connection to handle requests
//
// port number eg: "8080"
func (s *Server) Start(port string) {
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
		go s.HandleConnection(conn)
	}
}

func (s *Server) readOpCode(reader *bufio.Reader) (byte, error) {
	return reader.ReadByte()
}
func (s *Server) handleError(conn net.Conn, errStatus byte, message string) {
	sendResponse(conn, errStatus, []byte(message))
}
func (s *Server) handleSetOperatin(reader *bufio.Reader, conn net.Conn) {
	keyLengthBuf := make([]byte, 2)
	_, err := reader.Read(keyLengthBuf)
	if err != nil {
		s.handleError(conn, protocol.StatusErrReadKey, fmt.Sprintf("Error reading key length: %v", err))
		log.Println("Error reading key length:", err)
		return
	}
	keyLength := binary.BigEndian.Uint16(keyLengthBuf) // converting []byte of length 2 into unsigned integer
	// eg: converts [0x00, 0x03] -> 3
	keyBuf := make([]byte, keyLength)
	_, err = reader.Read(keyBuf)

	if err != nil {
		s.handleError(conn, protocol.StatusErrReadKey, fmt.Sprintf("Error reading key: %v", err))
		log.Println("Error reading key:", err)
		return
	}

	valueLengthBuf := make([]byte, 2)
	_, err = reader.Read(valueLengthBuf)
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
	s.dbService.Set(key, value)
	sendResponse(conn, 0x00, []byte(STATUS_SET_SUCCESS))
	log.Printf("SET operation for key: %s, value: %s\n", keyBuf, valueBuf)
}

// func handleConnection(conn net.Conn) {
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		opCode, err := s.readOpCode(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				return
			}
			log.Println("Error reading opCode:", err)
			s.handleError(conn, protocol.StatusErrReadOp, fmt.Sprintf("Error reading operation code %v", err))
			return
		}
		switch opCode {
		// Only process value for SET operation (0x01)
		case protocol.OpSet:
			s.handleSetOperatin(reader, conn)
		case protocol.OpGet:
			// handle the GET operation
			log.Println("GET operation ")
		}

	}
}
