package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/cirias/remote-controller/pkg/ac"
	"github.com/pkg/errors"
)

var serverAddr = flag.String("server", "", "address of the IR code server")

/*
 * var onCode = [8]byte{0x29, 0x09, 0x00, 0x50, 0x00, 0x00, 0x00, 0xC0}
 * var offCode = [8]byte{0x21, 0x09, 0x00, 0x50, 0x00, 0x00, 0x00, 0x40}
 */

func main() {
	flag.Parse()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(w, req)
			return
		}

		conn, err := net.Dial("tcp", *serverAddr)
		if err != nil {
			http.Error(w, errors.Wrapf(err, "could not dial %s", *serverAddr).Error(), 500)
			return
		}
		defer conn.Close()

		var s ac.State
		if err := json.NewDecoder(req.Body).Decode(&s); err != nil {
			http.Error(w, errors.Wrapf(err, "could not decode state").Error(), 400)
			return
		}

		code := s.Code()
		log.Printf("sending code: %X", code)
		if err := ac.Send(conn, code); err != nil {
			http.Error(w, errors.Wrapf(err, "could not send code").Error(), 500)
			return
		}

		io.WriteString(w, "ok")
	}

	http.HandleFunc("/state", helloHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))

	/*
	 * s := &ac.State{
	 *   Power: false,
	 *   Mode:  "cool",
	 *   Wind:  2,
	 *   Temp:  30,
	 * }
	 * fmt.Printf("sending code: %X\n", s.Code())
	 * if err := ac.Send(conn, s.Code()); err != nil {
	 *   log.Fatalln("could not send code:", err)
	 * }
	 */
}
