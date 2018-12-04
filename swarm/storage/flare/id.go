package flare

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

const proofParameterLength = 4
const proofTypeLength = 1
const proofLength = 32
const IDLength = storage.AddressLength
const TopicLength = IDLength - proofParameterLength - proofTypeLength

type Topic [TopicLength]byte

type ProofType byte

const DefaultProofType = ProofType(0)

type ID struct {
	ProofType      ProofType
	ProofParameter [4]byte
	Topic          Topic
}

var ErrIncorrectIDLength = errors.New("Incorrect ID Length.")

func (id *ID) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, IDLength)
	data = append(data, byte(id.ProofType))
	data = append(data, id.ProofParameter[:]...)
	data = append(data, id.Topic[:]...)
	return data, nil
}

func (id *ID) UnmarshalBinary(data []byte) error {
	if len(data) != IDLength {
		return ErrIncorrectIDLength
	}
	cursor := 0
	id.ProofType = ProofType(data[cursor])
	cursor += proofTypeLength
	copy(id.ProofParameter[:], data[cursor:cursor+proofParameterLength])
	cursor += proofParameterLength
	copy(id.Topic[:], data[cursor:cursor+TopicLength])
	cursor += TopicLength
	return nil
}

type Nonce struct {
	n1 uint64
	n2 uint64
}

type Value struct {
	Proof common.Hash
	Nonce Nonce
	Data  []byte
}

func (v *Value) Verify(id *ID) bool {

	return false
}

func (v *Value) MarshalBinary() ([]byte, error) {

	return nil, nil
}

func (v *Value) UnmarshalBinary(data []byte) error {

	return nil
}
