package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

type dumper struct {
	w      io.Writer
	label  string
	dumpTo io.Writer
}

func (d *dumper) Write(b []byte) (int, error) {
	message := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), d.label)
	io.WriteString(d.dumpTo, message)
	io.WriteString(d.dumpTo, hex.Dump(b))
	return d.w.Write(b)
}
