package flare

import (
	"encoding/binary"
	"errors"
)

type Nonce struct {
	N0 uint64
	N1 uint64
}

func (n *Nonce) Inc() {
	n.N0++
	if n.N0 == 0 {
		n.N1++
	}
}

func (n *Nonce) BinaryPut(data []byte) error {
	if len(data) != NonceLength {
		return errors.New("Incorrect Nonce length")
	}
	binary.LittleEndian.PutUint64(data[:8], n.N0)
	binary.LittleEndian.PutUint64(data[8:16], n.N1)
	return nil
}

func (n *Nonce) BinaryLength() int {
	return NonceLength
}
func (n *Nonce) BinaryGet(data []byte) error {
	if len(data) != NonceLength {
		return errors.New("Incorrect Nonce length")
	}
	n.N0 = binary.LittleEndian.Uint64(data[:8])
	n.N1 = binary.LittleEndian.Uint64(data[8:16])
	return nil
}
