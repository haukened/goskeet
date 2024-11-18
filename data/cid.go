package data

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

// CBOR_NULL represents the CBOR (Concise Binary Object Representation) encoding for the null value.
const (
	CBOR_NULL = 0xf6
)

// This is the "cid-link" type from the atproto data model.
type CIDLink struct {
	link cid.Cid
}

// jsonLink represents a structure that holds a single JSON field named "$link".
// The Link field is a string that maps to the "$link" key in the JSON representation.
type jsonLink struct {
	Link string `json:"$link"`
}

// String returns the string representation of the CIDLink.
// It delegates the call to the String method of the underlying link.
func (c CIDLink) String() string {
	return c.link.String()
}

// MarshalJSON implements the json.Marshaler interface for the CIDLink type.
// It returns the JSON encoding of the CIDLink. If the CIDLink is not defined,
// it returns an error indicating that a nil CIDLink was attempted to be marshaled.
func (c CIDLink) MarshalJSON() ([]byte, error) {
	if !c.link.Defined() {
		return nil, fmt.Errorf("tried to marshal nil cid-link")
	}
	return json.Marshal(jsonLink{Link: c.link.String()})
}

// UnmarshalJSON implements the json.Unmarshaler interface for the CIDLink type.
// It decodes a JSON-encoded byte slice into a CIDLink object.
// The function first unmarshals the data into a temporary jsonLink struct.
// Then, it decodes the CID from the jsonLink and assigns it to the CIDLink.
// If any error occurs during unmarshalling or decoding, it returns the error.
func (c *CIDLink) UnmarshalJSON(data []byte) error {
	var jl jsonLink
	if err := json.Unmarshal(data, &jl); err != nil {
		return err
	}
	link, err := cid.Decode(jl.Link)
	if err != nil {
		return err
	}
	c.link = link
	return nil
}

// MarshalCBOR encodes the CIDLink into CBOR format and writes it to the provided io.Writer.
// If the CIDLink is nil, it writes a CBOR null value. If the CIDLink is not defined, it returns an error.
// It uses CTAP2 encoding options for CBOR encoding.
func (c *CIDLink) MarshalCBOR(w io.Writer) error {
	enc, err := cbor.CTAP2EncOptions().EncMode()
	if err != nil {
		return err
	}
	if c == nil {
		_, err := w.Write([]byte{CBOR_NULL})
		return err
	}
	if !c.link.Defined() {
		return fmt.Errorf("tried to marshal nil cid-link")
	}
	return enc.NewEncoder(w).Encode(c.link)
}

// UnmarshalCBOR decodes CBOR data from the provided io.Reader and
// stores the result in the CIDLink receiver. It uses a CBOR decoder
// to read the data and populate the link field of the CIDLink.
// Returns an error if the decoding process fails.
func (c *CIDLink) UnmarshalCBOR(r io.Reader) error {
	dec := cbor.NewDecoder(r)
	return dec.Decode(&c.link)
}
