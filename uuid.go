// gouuid package implements UUID (version 4, random) type and methods for the manipulation of it.
// It implements the encoding/json.Marshaler and encoding/json.Unmarshaler interface to allow easy
// handling with JSON.
package gouuid
import (
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

const (
	UUIDLen = 16
)

type UUID [UUIDLen]byte

// New generates and returns new UUID v4 (generated randomly).
func New() (u UUID) {
	if _, e := io.ReadFull(rand.Reader, u[:]); e != nil {
		panic("error reading from random source: " + e.Error())
	}
	u[6] = u[6]>>4 | 0x40 // set version number
	u[8] = u[8] & ^uint8(1 << 6)
	u[8] |= 1 << 7
	return
}

func parseShortString(s string) (u UUID, e error) {
	b := []byte(s)
	if hex.DecodedLen(len(s)) != UUIDLen {
		e = fmt.Errorf("uuid: wrong string length for decode")
		return
	}
	_, e = hex.Decode(u[:], b)
	return
}

// ParseString converts a string (hex uuid, can include dashes) to UUID.
func ParseString(s string) (UUID, error) {
	s = strings.Replace(s, "-", "", -1)
	return parseShortString(s)
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// ShortString returns short string representation (without dashes) of UUID.
// Example: b7c016dc2ba4a68db368a97da9f43cee
func (u UUID) ShortString() string {
	return fmt.Sprintf("%x", u[:])
}

func (u UUID) Equal(a UUID) bool {
	for i, v := range u {
		if v != a[i] {
			return false
		}
	}
	return true
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.ShortString() + "\""), nil
}

func (u *UUID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 {
		return fmt.Errorf("uuid: JSON value is too short for UUID")
	}
	x, e := ParseString(string(b[1 : len(b)-1]))
	copy((*u)[:], x[:])
	return e
}
