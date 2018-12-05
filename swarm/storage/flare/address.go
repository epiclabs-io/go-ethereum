package flare

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

const proofParameterLength = 4
const proofTypeLength = 1
const AddressLength = storage.AddressLength
const TopicLength = AddressLength - proofParameterLength - proofTypeLength
const NonceLength = 8 + 8
const nonceOffset = AddressLength
const ValueHeaderLength = AddressLength + NonceLength

type Topic [TopicLength]byte

type ProofType byte

const DefaultProofType = ProofType(0)

type Address struct {
	ProofType      ProofType
	ProofParameter [4]byte
	Topic          Topic
}

var ErrIncorrectIDLength = errors.New("Incorrect ID Length.")

func (addr *Address) BinaryPut(data []byte) error {
	if len(data) != AddressLength {
		return errors.New("Incorrect address length")
	}
	cursor := 0

	data[cursor] = byte(addr.ProofType)
	cursor += proofTypeLength

	copy(data[cursor:cursor+proofParameterLength], addr.ProofParameter[:])
	cursor += proofParameterLength

	copy(data[cursor:cursor+TopicLength], addr.Topic[:])
	cursor += TopicLength

	return nil
}

func (addr *Address) BinaryLength() int {
	return AddressLength
}

func (addr *Address) BinaryGet(data []byte) error {
	if len(data) != AddressLength {
		return errors.New("Incorrect address length")
	}
	cursor := 0

	addr.ProofType = ProofType(data[cursor])
	cursor += proofTypeLength

	copy(addr.ProofParameter[:], data[cursor:cursor+proofParameterLength])
	cursor += proofParameterLength

	copy(addr.Topic[:], data[cursor:cursor+TopicLength])
	cursor += TopicLength
	return nil
}

func (addr *Address) Bytes() []byte {
	data := make([]byte, AddressLength)
	addr.BinaryPut(data)
	return data
}

func (addr *Address) Target() []byte {
	//https://bitcoin.org/en/developer-reference#target-nbits
	target := make([]byte, common.HashLength)
	exponent := int(addr.ProofParameter[0]) - 3
	copy(addr.ProofParameter[1:])
	return nil
}
