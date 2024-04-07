package protocol

const (
	OpSet            byte = 0x01
	OpGet            byte = 0x02
	StatusSuccess    byte = 0x10 // General success status for operations
	StatusSetSuccess byte = 0x11 // Success specifically for SET operations
	StatusGetSuccess byte = 0x12 // Success specifically for GET operations
	StatusErrReadKey byte = 0xFE // Error when reading the key
	StatusErrReadOp  byte = 0xFF // Error when reading the operation code
)

type Request struct {
	Operation byte
	Key       []byte
	Value     []byte
}

type Response struct {
	Status byte
	Data   []byte
}
