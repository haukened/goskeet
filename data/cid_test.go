package data

import (
	"bytes"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

// This is a random CID grabbed from the Firehose. Sorry if its yours.
const (
	TEST_CID     = "bafyreibfd77vb2setujncomtz3j6xswrmiuxlykora6nogxbr4arhqu2ye"
	INVALID_CID  = "bafyreibfd77vb2setujncomtz3j6xswrmiuxlykora6nogxbr4arhqu2y" // too short
	TEST_JSON    = `{"$link": "bafyreibfd77vb2setujncomtz3j6xswrmiuxlykora6nogxbr4arhqu2ye"}`
	INVALID_JSON = `{$link: "bafyreibfd77vb2setujncomtz3j6xswrmiuxlykora6nogxbr4arhqu2ye"}` // missing quotes
)

func TestCIDLink_String(t *testing.T) {
	c, err := cid.Decode(TEST_CID)
	assert.NoError(t, err)
	link := CIDLink{link: c}
	assert.Equal(t, c.String(), link.String())
}

func TestCIDLinkMarshalJSONNull(t *testing.T) {
	link := CIDLink{}
	_, err := link.MarshalJSON()
	assert.Error(t, err)
}

func TestCIDLinkMarshalInvalidJSON(t *testing.T) {
	var link CIDLink
	err := link.UnmarshalJSON([]byte(INVALID_JSON))
	assert.Error(t, err)
}

func TestCIDLinkMarshalInvalidCID(t *testing.T) {
	var link CIDLink
	err := link.UnmarshalJSON([]byte(`{"$link": "` + INVALID_CID + `"}`))
	assert.Error(t, err)
}

func TestMarshalNilCBOR(t *testing.T) {
	link := &CIDLink{}
	link = nil
	var buf bytes.Buffer
	err := link.MarshalCBOR(&buf)
	assert.NoError(t, err)
	// CBOR NULL is 0xf6
	assert.Equal(t, []byte{0xf6}, buf.Bytes())
}

func TestCIDLinkMarshalCBORNull(t *testing.T) {
	link := CIDLink{}
	var buf bytes.Buffer
	err := link.MarshalCBOR(&buf)
	assert.Error(t, err)
}

func TestCIDLink_MarshalJSON(t *testing.T) {
	c, err := cid.Decode(TEST_CID)
	assert.NoError(t, err)
	link := CIDLink{link: c}
	data, err := link.MarshalJSON()
	assert.NoError(t, err)
	expected := `{"$link":"` + c.String() + `"}`
	assert.JSONEq(t, expected, string(data))
}

func TestCIDLink_UnmarshalJSON(t *testing.T) {
	var link CIDLink
	err := link.UnmarshalJSON([]byte(TEST_JSON))
	assert.NoError(t, err)
	expected, err := cid.Decode(TEST_CID)
	assert.NoError(t, err)
	assert.Equal(t, expected, link.link)
}

func TestCIDLink_MarshalCBOR(t *testing.T) {
	c, err := cid.Decode(TEST_CID)
	assert.NoError(t, err)
	link := CIDLink{link: c}
	var buf bytes.Buffer
	err = link.MarshalCBOR(&buf)
	assert.NoError(t, err)
	dec := cbor.NewDecoder(&buf)
	var decoded cid.Cid
	err = dec.Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, c, decoded)
}

func TestCIDLink_UnmarshalCBOR(t *testing.T) {
	c, err := cid.Decode(TEST_CID)
	assert.NoError(t, err)
	var buf bytes.Buffer
	enc, err := cbor.CTAP2EncOptions().EncMode()
	assert.NoError(t, err)
	err = enc.NewEncoder(&buf).Encode(c)
	assert.NoError(t, err)
	var link CIDLink
	err = link.UnmarshalCBOR(&buf)
	assert.NoError(t, err)
	assert.Equal(t, c, link.link)
}
