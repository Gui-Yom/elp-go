package main

import (
	"bufio"
	"elp-go/internal"
	"github.com/fxamacker/cbor/v2"
	"net"
	"reflect"
)

type Remote struct {
	conn    *net.TCPConn
	bufw    *bufio.Writer
	encoder *cbor.Encoder
	decoder *cbor.Decoder
}

func NewRemote(conn *net.TCPConn) *Remote {
	tags := cbor.NewTagSet()
	tags.Add(cbor.TagOptions{EncTag: cbor.EncTagRequired, DecTag: cbor.DecTagRequired}, reflect.TypeOf(internal.MoveTask{}), 279)
	encMode, _ := cbor.PreferredUnsortedEncOptions().EncModeWithTags(tags)
	decMode, _ := cbor.DecOptions{}.DecModeWithTags(tags)

	// Buffer sizes are arbitrary
	bufw := bufio.NewWriterSize(conn, 4096)
	return &Remote{
		conn:    conn,
		bufw:    bufw,
		encoder: encMode.NewEncoder(bufw),
		decoder: decMode.NewDecoder(bufio.NewReaderSize(conn, 4096)),
	}
}

func (rem *Remote) Send(v interface{}) error {
	err := rem.encoder.Encode(v)
	// With a call to Send, we expect the data to be sent immediately
	rem.bufw.Flush()
	return err
}

// Recv v must be a pointer
func (rem *Remote) Recv(v interface{}) error {
	return rem.decoder.Decode(v)
}

func (rem *Remote) Close() error {
	return rem.conn.Close()
}
