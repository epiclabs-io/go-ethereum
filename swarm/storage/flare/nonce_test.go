package flare_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage/flare"
	"github.com/ethereum/go-ethereum/swarm/storage/flare/binser"
)

func getTestNonce() *flare.Nonce {
	return &flare.Nonce{
		N0: 2352341634,
		N1: 8679432341257,
	}
}

func TestNonceSerializerDeserializer(t *testing.T) {
	binser.TestBinarySerializerRecovery(t, getTestNonce(), "0x82e2358c0000000009477bd6e4070000")
}

func TestNonceSerializerLengthCheck(t *testing.T) {
	binser.TestBinarySerializerLengthCheck(t, getTestNonce())
}
