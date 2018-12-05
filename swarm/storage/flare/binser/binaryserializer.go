// TODO: move to a common package to share with feed package

package binser

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type BinarySerializer interface {
	BinaryPut(serializedData []byte) error
	BinaryLength() int
	BinaryGet(serializedData []byte) error
}

// Values interface represents a string key-value store
// useful for building query strings
type Values interface {
	Get(key string) string
	Set(key, value string)
}

type valueSerializer interface {
	FromValues(values Values) error
	AppendValues(values Values)
}

// Hex serializes the structure and converts it to a hex string
func Hex(bin BinarySerializer) string {
	b := make([]byte, bin.BinaryLength())
	bin.BinaryPut(b)
	return hexutil.Encode(b)
}

// KV mocks a key value store
type KV map[string]string

func (kv KV) Get(key string) string {
	return kv[key]
}
func (kv KV) Set(key, value string) {
	kv[key] = value
}

func CompareByteSliceToExpectedHex(t *testing.T, variableName string, actualValue []byte, expectedHex string) {
	if hexutil.Encode(actualValue) != expectedHex {
		t.Fatalf("%s: Expected %s to be %s, got %s", t.Name(), variableName, expectedHex, hexutil.Encode(actualValue))
	}
}

func TestBinarySerializerRecovery(t *testing.T, bin BinarySerializer, expectedHex string) {
	name := reflect.TypeOf(bin).Elem().Name()
	serialized := make([]byte, bin.BinaryLength())
	if err := bin.BinaryPut(serialized); err != nil {
		t.Fatalf("%s.BinaryPut error when trying to serialize structure: %s", name, err)
	}

	CompareByteSliceToExpectedHex(t, name, serialized, expectedHex)

	recovered := reflect.New(reflect.TypeOf(bin).Elem()).Interface().(BinarySerializer)
	if err := recovered.BinaryGet(serialized); err != nil {
		t.Fatalf("%s.BinaryGet error when trying to deserialize structure: %s", name, err)
	}

	if !reflect.DeepEqual(bin, recovered) {
		t.Fatalf("Expected that the recovered %s equals the marshalled %s", name, name)
	}

	serializedWrongLength := make([]byte, 1)
	copy(serializedWrongLength[:], serialized)
	if err := recovered.BinaryGet(serializedWrongLength); err == nil {
		t.Fatalf("Expected %s.BinaryGet to fail since data is too small", name)
	}
}

func TestBinarySerializerLengthCheck(t *testing.T, bin BinarySerializer) {
	name := reflect.TypeOf(bin).Elem().Name()
	// make a slice that is too small to contain the metadata
	serialized := make([]byte, bin.BinaryLength()-1)

	if err := bin.BinaryPut(serialized); err == nil {
		t.Fatalf("Expected %s.BinaryPut to fail, since target slice is too small", name)
	}
}

func TestValueSerializer(t *testing.T, v valueSerializer, expected KV) {
	name := reflect.TypeOf(v).Elem().Name()
	kv := make(KV)

	v.AppendValues(kv)
	if !reflect.DeepEqual(expected, kv) {
		expj, _ := json.Marshal(expected)
		gotj, _ := json.Marshal(kv)
		t.Fatalf("Expected %s.AppendValues to return %s, got %s", name, string(expj), string(gotj))
	}

	recovered := reflect.New(reflect.TypeOf(v).Elem()).Interface().(valueSerializer)
	err := recovered.FromValues(kv)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(recovered, v) {
		t.Fatalf("Expected recovered %s to be the same", name)
	}
}
