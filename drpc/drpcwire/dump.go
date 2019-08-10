// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package drpcwire

import (
	"fmt"
	"io"

	"storj.io/storj/drpc"
)

type Dumper struct {
	out io.Writer
	err error
	buf []byte
}

func NewDumper(out io.Writer) *Dumper {
	return &Dumper{out: out}
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}
	d.buf = append(d.buf, p...)

	fmt.Fprintf(d.out, "write: %x\n", p)

	defer func() {
		if err != nil && d.err == nil {
			d.err = err
		}
	}()

	for {
		advance, token, err := PacketScanner(d.buf, false)
		if err != nil {
			return len(p), err
		} else if token == nil {
			return len(p), nil
		}

		rem, pkt, ok, err := ParsePacket(token)
		if !ok || err != nil || len(rem) > 0 {
			return len(p), drpc.InternalError.New("invalid parse after scanner")
		}
		d.buf = d.buf[advance:]

		if _, err := fmt.Fprintf(
			d.out, "     | pid:<%d,%d> kind:%d cont:%-5v start:%-5v len:%-4d data:%x\n",
			pkt.StreamID, pkt.MessageID,
			pkt.PayloadKind,
			pkt.Continuation,
			pkt.Starting,
			pkt.Length,
			pkt.Data,
		); err != nil {
			return len(p), err
		}
	}
}
