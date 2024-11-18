package data

import (
	"bytes"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

const TEST_BYTES = "dGVzdCBkYXRh"
const TEST_BYTES_JSON = `{"$bytes":"dGVzdCBkYXRh"}`
const TEST_BYTES_STRING = "test data"

func TestBytes_MarshalJSON(t *testing.T) {
	b := Bytes(TEST_BYTES_STRING)
	data, err := b.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, TEST_BYTES_JSON, string(data))
}

func TestBytes_UnmarshalJSON(t *testing.T) {
	var b Bytes
	err := b.UnmarshalJSON([]byte(TEST_BYTES_JSON))
	assert.NoError(t, err)
	assert.Equal(t, Bytes(TEST_BYTES), b)
}

func TestBytes_MarshalJSON_Nil(t *testing.T) {
	var b Bytes
	data, err := b.MarshalJSON()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestBytes_MarshalCBOR(t *testing.T) {
	b := Bytes(TEST_BYTES_STRING)
	var buf bytes.Buffer
	err := b.MarshalCBOR(&buf)
	assert.NoError(t, err)
	dec := cbor.NewDecoder(&buf)
	var decoded []byte
	err = dec.Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), decoded)
}

func TestBytes_UnmarshalCBOR(t *testing.T) {
	b := Bytes(TEST_BYTES_STRING)
	var buf bytes.Buffer
	enc, err := cbor.CTAP2EncOptions().EncMode()
	assert.NoError(t, err)
	err = enc.NewEncoder(&buf).Encode(b.Bytes())
	assert.NoError(t, err)
	var decoded Bytes
	err = decoded.UnmarshalCBOR(&buf)
	assert.NoError(t, err)
	assert.Equal(t, b, decoded)
}

func TestBytes_MarshalCBOR_Nil(t *testing.T) {
	var b *Bytes
	var buf bytes.Buffer
	err := b.MarshalCBOR(&buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xf6}, buf.Bytes()) // CBOR NULL is 0xf6
}
