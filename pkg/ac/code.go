package ac

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/pkg/errors"
)

var kEndian = binary.LittleEndian

const kGreeHdrMark uint16 = 9000
const kGreeHdrSpace uint16 = 4000
const kGreeBitMark uint16 = 620
const kGreeOneSpace uint16 = 1600
const kGreeZeroSpace uint16 = 540
const kGreeMsgSpace uint16 = 19000

const kGreeStateLength = 8
const kGreeBits = kGreeStateLength * 8

const kGreeBlockFooter byte = 2
const kGreeBlockFooterBits byte = 3

const kHeader = 2 // Usual nr. of header entries.
const kFooter = 2 // Usual nr. of footer (stop bits) entries.

func codeToSignal(code [8]byte) []uint16 {
	size := 2*(kGreeBits+kGreeBlockFooterBits) + kHeader + kFooter
	log.Println("calulated size:", size)
	signal := make([]uint16, 0, size)

	signal = append(signal, kGreeHdrMark)
	signal = append(signal, kGreeHdrSpace)

	for i := 0; i < 4; i++ {
		signal = appendBits(signal, code[i], 8)
	}

	signal = appendBits(signal, kGreeBlockFooter, kGreeBlockFooterBits)

	signal = append(signal, kGreeBitMark)
	signal = append(signal, kGreeMsgSpace)

	for i := 4; i < 8; i++ {
		signal = appendBits(signal, code[i], 8)
	}

	signal = append(signal, kGreeBitMark)
	signal = append(signal, kGreeMsgSpace)
	log.Println("actual size:", len(signal))

	return signal
}

func appendBits(signal []uint16, b byte, nbits uint8) []uint16 {
	result := signal
	var bit uint8
	for bit = 0; bit < nbits; bit++ {
		if ((b >> bit) & 1) == 1 { // Send a 1
			result = append(result, kGreeBitMark)
			result = append(result, kGreeOneSpace)
		} else { // Send a 0
			result = append(result, kGreeBitMark)
			result = append(result, kGreeZeroSpace)
		}
	}
	return result
}

func Send(r io.Writer, code [8]byte) error {
	signal := codeToSignal(code)

	size := uint16(binary.Size(signal))
	if err := binary.Write(r, kEndian, size); err != nil {
		return errors.Wrapf(err, "could not send size")
	}

	if err := binary.Write(r, kEndian, signal); err != nil {
		return errors.Wrapf(err, "could not send signal")
	}

	return nil
}
