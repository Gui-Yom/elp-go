package shared

import (
	"bufio"
	"github.com/fxamacker/cbor/v2"
	"net"
)

type Remote struct {
	conn    *net.TCPConn
	encoder *cbor.Encoder
	decoder *cbor.Decoder
}

func NewRemote(conn *net.TCPConn) *Remote {
	// Buffer sizes are arbitrary
	return &Remote{
		conn:    conn,
		encoder: cbor.NewEncoder(bufio.NewWriterSize(conn, 4096)),
		decoder: cbor.NewDecoder(bufio.NewReaderSize(conn, 4096)),
	}
}

func (rem *Remote) Send(v interface{}) {
	rem.encoder.Encode(v)
}

// Recv v must be a pointer
func (rem *Remote) Recv(v interface{}) {
	rem.decoder.Decode(v)
}

func (rem *Remote) Close() {
	rem.conn.Close()
}
