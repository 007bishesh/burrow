package crypto

import (
	"bytes"
	"fmt"
	"io"
)

type CurveType int8

const (
	CurveTypeSecp256k1 CurveType = iota
	CurveTypeEd25519
)

func (k CurveType) String() string {
	switch k {
	case CurveTypeSecp256k1:
		return "secp256k1"
	case CurveTypeEd25519:
		return "ed25519"
	default:
		return "unknown"
	}
}

// Get this CurveType's 8 bit identifier as a byte
func (k CurveType) Byte() byte {
	return byte(k)
}

func CurveTypeFromString(s string) (CurveType, error) {
	switch s {
	case "secp256k1":
		return CurveTypeSecp256k1, nil
	case "ed25519":
		return CurveTypeEd25519, nil
	default:
		var k CurveType
		return k, ErrInvalidCurve(s)
	}
}

type ErrInvalidCurve string

func (err ErrInvalidCurve) Error() string {
	return fmt.Sprintf("invalid curve type")
}

// The types in this file allow us to control serialisation of keys and signatures, as well as the interface
// exposed regardless of crypto library

type Signer interface {
	Sign(msg []byte) (Signature, error)
}

// Signable is an interface for all signable things.
// It typically removes signatures before serializing.
type Signable interface {
	WriteSignBytes(chainID string, w io.Writer, n *int, err *error)
}

// SignBytes is a convenience method for getting the bytes to sign of a Signable.
func SignBytes(chainID string, o Signable) []byte {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	o.WriteSignBytes(chainID, buf, n, err)
	if *err != nil {
		panic(fmt.Sprintf("could not write sign bytes for a signable: %s", *err))
	}
	return buf.Bytes()
}
