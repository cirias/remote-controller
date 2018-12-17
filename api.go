package main

import (
	"encoding/binary"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type remoteController struct {
	code []uint16
	conn net.Conn
}

func (c *remoteController) send() error {
	size := uint16(binary.Size(c.code))
	if err := binary.Write(c.conn, kEndian, size); err != nil {
		return errors.Wrapf(err, "could not send size")
	}

	if err := binary.Write(c.conn, kEndian, c.code); err != nil {
		return errors.Wrapf(err, "could not send code")
	}

	return nil
}

func (c *remoteController) setPower(on bool) {

}

func (c *remoteController) setTemp(t int) {
	c.code = []uint16{}
}

func (c *remoteController) setWindSpeed(s int) {

}

func (c *remoteController) setMode(m string) {

}

func (c *remoteController) Code() []uint16 {
	code := make([]uint16, len(c.code))
	copy(code, c.code)
	return code
}

func handleIR(w http.ResponseWriter, req *http.Request) {

}

func newApiHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ir", handleIR)

	return mux
}
