package flare

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

type Value struct {
	Addr  Address
	Nonce Nonce
	Data  []byte
}

func (v *Value) Verify(chunkAddr storage.Address) bool {
	if !bytes.Equal(chunkAddr[:], v.Addr.Bytes()) {
		return false
	}

	data, err := v.MarshalBinary()
	if err != nil {
		return false
	}
	target := v.Addr.Target()
	hasher := sha3.NewKeccak256()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	if bytes.Compare(target, hash) > 0 {
		return true
	}
	return false
}

func (v *Value) PoW() error {
	target := v.Addr.Target()
	// set proof & nonce to 0
	v.Nonce = Nonce{}
	hasher := sha3.NewKeccak256()
	data, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	for {
		hasher.Reset()
		hasher.Write(data)
		hash := hasher.Sum(nil)
		if bytes.Compare(target, hash) > 0 {
			return nil
		}
		v.Nonce.Inc()
		v.Nonce.BinaryPut(data[nonceOffset : nonceOffset+NonceLength])
	}
}

func (v *Value) BinaryPut(data []byte) error {
	if len(data) != v.BinaryLength() {
		return errors.New("Incorrect Value length")
	}
	cursor := 0

	if err := v.Addr.BinaryPut(data[cursor : cursor+AddressLength]); err != nil {
		return err
	}
	cursor += AddressLength

	if err := v.Nonce.BinaryPut(data[cursor : cursor+NonceLength]); err != nil {
		return err
	}
	cursor += NonceLength

	copy(data[cursor:], v.Data)

	return nil
}

func (v *Value) BinaryLength() int {
	return ValueHeaderLength + len(v.Data)
}
func (v *Value) BinaryGet(data []byte) error {
	if len(data) < ValueHeaderLength {
		return errors.New("Incorrect Value length")
	}
	cursor := 0

	if err := v.Addr.BinaryGet(data[cursor : cursor+AddressLength]); err != nil {
		return err
	}
	cursor += AddressLength

	if err := v.Nonce.BinaryGet(data[cursor : cursor+NonceLength]); err != nil {
		return err
	}
	cursor += NonceLength

	v.Data = make([]byte, len(data)-cursor)
	copy(v.Data[:], data[cursor:])
	return nil
}

func (v *Value) MarshalBinary() ([]byte, error) {
	data := make([]byte, v.BinaryLength())
	return data, v.BinaryPut(data)
}

func (v *Value) UnmarshalBinary(data []byte) error {
	return v.BinaryGet(data)
}

func (v *Value) toChunk() (storage.Chunk, error) {
	valueBytes, err := v.MarshalBinary()
	if err != nil {
		return nil, err
	}
	chunkAddr := make(storage.Address, storage.AddressLength)
	v.Addr.BinaryPut([]byte(chunkAddr))

	return storage.NewChunk(chunkAddr, valueBytes), nil
}

func (v *Value) fromChunk(chunk storage.Chunk) error {
	return v.UnmarshalBinary(chunk.Data())
}
