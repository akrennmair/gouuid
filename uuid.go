// Copyright 2011 Dmitry Chestnykh

// UUID package implements UUID (version 4, random) type and methods for the manipulation of it.
package uuid

import (
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
)

const (
	UUIDLen = 16
)

type UUID [UUIDLen]byte

// New generates and returns new UUID v4 (generated randomly).
func New() (u UUID) {
	rand.Read(u[:])
	u[6] = u[6]>>4 | 0x40 // set version number
	u[8] &^= 1 << 6       // set 6th bit to 0
	u[8] |= 1 << 7        // set 7th bit to 1
	return
}

// NewShortString converts a short string (hex uuid without dashes) to UUID.
func NewShortString(s string) (u UUID, err os.Error) {
	b := []byte(s)
	if hex.DecodedLen(len(s)) != UUIDLen {
		err = os.NewError("uuid: wrong string length for decode")
		return
	}
	_, err = hex.Decode(u[:], b)
	return
}

// NewShortString converts a string (hex uuid, can include dashes) to UUID.
func NewString(s string) (UUID, os.Error) {
	s = strings.Replace(s, "-", "", -1)
	return NewShortString(s)
}

// String returns string representation of UUID.
// Example: b7c016dc-2ba4-a68d-b368-a97da9f43cee
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// ShortString returns short string representation (without dashes) of UUID.
// Example: b7c016dc2ba4a68db368a97da9f43cee
func (u UUID) ShortString() string {
	return fmt.Sprintf("%x", u[:])
}

// Bytes returns a byte slice of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Equal returns a boolean reporting whether UUID equals another one (a).
func (u UUID) Equal(a UUID) bool {
	for i, v := range u {
		if v != a[i] {
			return false
		}
	}
	return true
}

// MarshalJSON encodes UUID pointer into JSON representation.
func (u *UUID) MarshalJSON() ([]byte, os.Error) {
	return []byte("\"" + u.ShortString() + "\""), nil
}

// UnmarshalJSON decodes UUID pointer from JSON representation.
func (u *UUID) UnmarshalJSON(b []byte) os.Error {
	if len(b) < 3 {
		return os.NewError("uuid: JSON value is too short for UUID")
	}
	x, err := NewShortString(string(b[1 : len(b)-1]))
	copy((*u)[:], x[:])
	return err
}
