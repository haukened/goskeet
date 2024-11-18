package data

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
)

// represents the "bytes" type from the atproto data model
// in JSON, marshals to an object with $bytes key and base64-encoded data
type Bytes []byte

// JSONBytes represents a structure that holds a JSON-encoded string.
// The JSON representation uses the key "$bytes" to store the string value.
type JSONBytes struct {
	Bytes string `json:"$bytes"`
}

// Bytes converts the Bytes type to a slice of bytes ([]byte).
// It returns the underlying byte slice representation of the Bytes type.
func (b Bytes) Bytes() []byte {
	return []byte(b)
}

// MarshalJSON implements the json.Marshaler interface for the Bytes type.
// It encodes the Bytes value as a base64 string and wraps it in a JSONBytes struct.
// If the Bytes value is nil, it returns an error indicating that nil bytes cannot be marshaled.
// Returns the JSON-encoded byte slice or an error if the marshaling fails.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if b == nil {
		return nil, fmt.Errorf("cannot marshal nil $bytes")
	}
	jb := JSONBytes{
		Bytes: base64.RawStdEncoding.EncodeToString(b.Bytes()),
	}
	return json.Marshal(jb)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Bytes type.
// It decodes a JSON-encoded byte slice, which is expected to be a base64-encoded string,
// and assigns the decoded byte slice to the Bytes receiver.
//
// Parameters:
// - data: A byte slice containing the JSON-encoded data.
//
// Returns:
// - error: An error if the JSON unmarshalling or base64 decoding fails, otherwise nil.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	var jb JSONBytes
	if err := json.Unmarshal(data, &jb); err != nil {
		return err
	}
	decoded, err := base64.RawStdEncoding.DecodeString(jb.Bytes)
	if err != nil {
		return err
	}
	*b = Bytes(decoded)
	return nil
}

// MarshalCBOR encodes the Bytes object into CBOR format and writes it to the provided io.Writer.
// It uses CTAP2 encoding options for CBOR encoding.
// If the Bytes object is nil, it writes a CBOR null value to the writer.
// Returns an error if encoding fails or if there is an issue writing to the io.Writer.
func (b *Bytes) MarshalCBOR(w io.Writer) error {
	enc, err := cbor.CTAP2EncOptions().EncMode()
	if err != nil {
		return err
	}
	if b == nil {
		_, err := w.Write([]byte{CBOR_NULL})
		return err
	}
	return enc.NewEncoder(w).Encode(b.Bytes())
}

// UnmarshalCBOR decodes CBOR-encoded data from the provided io.Reader
// and stores the result in the Bytes struct. It returns an error if
// the decoding process fails.
func (b *Bytes) UnmarshalCBOR(r io.Reader) error {
	dec := cbor.NewDecoder(r)
	return dec.Decode(&b)
}
