// Package snö generates snowflake like 128 bit ids that uses milliseconds
// for the time element, a 7 byte ID, and a 20 bit secondary id.  A sne
// snoflake has a 11 bit sequence for uniqueness.  The sequence starts at a
// random position and increments with each Snowflake generated.  If there is
// a possibility of needing more than 1,024,000 Flakes per second for a given
// ID and secondary ID combination collisions will occur.
//
// Snö uses 41 bits for the time, representing the number of milliseconds
// since 1/1/2016 00:00:00.
//
// Data layout:
//    0-40  time, in milliseconds: 69 years
//   41-50  sequence
//   51-63  secondary ID
//  64-127  ID
package sno

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"
)

const (
	epoch        int64  = 1451606400 // 1/1/2016 00:00:00
	idBits              = 64
	sidBits             = 13
	sequenceBits        = 10
	sequenceMax  uint64 = 1<<sequenceBits - 1
	sidMask      int16  = -1 ^ (-1 << (sidBits % 16))
)

// Flake is the type for a 128 bit snowflake.
type Flake [16]byte

// Generator creates snowflakes for a given id
type Generator struct {
	id       []byte
	sid      int16
	sequence uint64
}

// New returns an initialized generator.  If the passed byte slice's length is
// greater than 8 bytes, the first 8 bytes will be used for the generator's id.
// If the passed byte slice's length is less than 8 bytes, the id will be left-
// padded with 0x00.  The id2 parameter is the secondary id: only the right-
// most 12 bits are used.  The generator's sequence is initialized with a
// random int within [0, 2^10).
func New(id []byte, id2 int16) Generator {
	var g Generator
	if len(id) < 8 {
		g.id = make([]byte, 8-len(id))
	}
	if len(id) > 8 {
		id = id[:8]
	}
	g.id = append(g.id, id...)
	g.sid = id2 & sidMask
	g.sequence = seed()
	return g
}

// Snowflake generates an Flake from the current time and next sequence value.
func (g *Generator) Snowflake() Flake {
	var flake Flake
	now := uint64(time.Now().UnixNano()/1000000 - epoch)
	v := atomic.AddUint64(&g.sequence, 1)
	v = v % sequenceMax
	flake[0] = byte(now >> 33)
	flake[1] = byte(now >> 25)
	flake[2] = byte(now >> 17)
	flake[3] = byte(now >> 9)
	flake[4] = byte(now >> 1)
	flake[5] = byte(now<<7) | uint8(v>>3)
	flake[6] = byte(v<<5) | byte(g.sid>>8)
	flake[7] = byte(g.sid)
	flake[8] = g.id[0]
	flake[9] = g.id[1]
	flake[10] = g.id[2]
	flake[11] = g.id[3]
	flake[12] = g.id[4]
	flake[13] = g.id[5]
	flake[14] = g.id[6]
	flake[15] = g.id[7]
	return flake
}

// ID returns the ID bytes associated with this generator.
func (g *Generator) ID() []byte {
	return g.id
}

// SID returns the secondary ID associated with this generator.
func (g *Generator) SID() int16 {
	return g.sid
}

// Time returns the Flake's timestamp as an int64.  The returned timestamp is
// milliseconds since Unix epoch.
func (f *Flake) Time() int64 {
	return int64(f[0])<<33 | int64(f[1])<<25 | int64(f[2])<<17 | int64(f[3])<<9 | int64(f[4])<<1 + epoch
}

// ID returns the Flake's ID as a []byte.  A snowflake's ID is the last 8
// bytes.  No assumptios are made about either the contents of those bytes
// or their layout, that is left up to the user.
func (f *Flake) ID() []byte {
	return f[8:]
}

// SID returns the Flake's secondary ID.  The secondary ID is 13 bits.
func (f *Flake) SID() int16 {
	return int16(f[6]<<3>>3)<<8 | int16(f[7])
}

func seed() uint64 {
	bi := big.NewInt(int64(sequenceMax) + 1)
	r, err := rand.Int(rand.Reader, bi)
	if err != nil {
		panic(fmt.Sprintf("entropy read error: %s", err))
	}
	return r.Uint64()
}
