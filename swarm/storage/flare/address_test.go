package flare_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage/flare"
	"github.com/ethereum/go-ethereum/swarm/storage/flare/binser"
)

func getTestAddress() *flare.Address {
	return &flare.Address{
		ProofType:      flare.DefaultProofType,
		ProofParameter: [4]byte{1, 2, 3, 4},
		Topic:          flare.Topic{},
	}
}

func TestAddressSerializeDeserialize(t *testing.T) {
	binser.TestBinarySerializerRecovery(t, getTestAddress(), "0x0001020304000000000000000000000000000000000000000000000000000000")
}

func TestAddressLengthCheck(t *testing.T) {
	binser.TestBinarySerializerLengthCheck(t, getTestAddress())
}
