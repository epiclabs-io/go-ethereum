package flare_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage/flare"
	"github.com/ethereum/go-ethereum/swarm/storage/flare/binser"
)

func getTestValue() *flare.Value {
	return &flare.Value{
		Addr:  *getTestAddress(),
		Nonce: *getTestNonce(),
		Data:  []byte("This is some data"),
	}
}

func TestValueSerializerDeserializer(t *testing.T) {
	binser.TestBinarySerializerRecovery(t, getTestValue(), "0x000102030400000000000000000000000000000000000000000000000000000082e2358c0000000009477bd6e40700005468697320697320736f6d652064617461	")
}

func TestValueSerializerLengthCheck(t *testing.T) {
	binser.TestBinarySerializerLengthCheck(t, getTestValue())
}

func TestPow(t *testing.T) {

}
