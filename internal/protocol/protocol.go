package protocol

const (
	OpSet byte = 0x01
	OpGet byte = 0x02
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
