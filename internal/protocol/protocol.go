package protocol

const (
	OpSet            byte = 0x01
	OpGet            byte = 0x02
	StatusSucces     byte = 0x00
	StatusErrReadOp  byte = 0xFF
	StatusSetSuccess byte = 0x01
	StatusErrReadKey byte = 0x0FE
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
