// Package sne generates snowflake like 128 bit ids that uses milliseconds
// for the time element. A sne Flake has a 41 bit time field, a 11 bit
// sequence, a 56 bit ID, and a 20 bit secondary id.  The time field holds
// milliseconds since 1/1/2016 00:00:00.  The sequence starts at a random
// position and increments with each Snowflake.  If there is the possibility
// of needing more than 2,048,000 Flakes per second for a given ID and
// secondary ID combination this package should not be used.
//
// Data layout:
//    0-40    time, in milliseconds: 69 years
//   41-51    sequence; for collision prevention
//   52-71    secondary id
//  72-127    primary ID
package sne

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"
)

const (
	epoch        int64  = 1451606400 // 1/1/2016 00:00:00
	idBits              = 56
	sidBits             = 20
	sequenceBits        = 11
	sequenceMax  uint64 = 1<<sequenceBits - 1
	sidMask      int32  = -1 ^ (-1 << (sidBits % 32))
)

// Flake is the type for a 128 bit snowflake.
type Flake [16]byte

// Generator creates snowflakes for a given id
type Generator struct {
	id       []byte
	sid      int32
	sequence uint64
}

// New returns an initialized generator.  If the passed byte slice's length is
// greater than 7 bytes, the first 7 bytes will be used for the generator's id.
// If the passed byte slice's length is less than 7 bytes, the id will be left-
// padded with 0x00.  The id2 parameter is the secondary id: only the right-
// most 20 bits are used.  The generator's sequence is initialized with a
// random int within [0, 2^11).
func New(id []byte, id2 int32) Generator {
	var g Generator
	if len(id) < 7 {
		g.id = make([]byte, 7-len(id))
	}
	if len(id) > 7 {
		id = id[:7]
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
	flake[5] = byte(now<<7) | uint8(v>>4)
	flake[6] = byte(v<<4) | byte(g.sid>>16)
	flake[7] = byte(g.sid >> 8)
	flake[8] = byte(g.sid)
	flake[9] = g.id[0]
	flake[10] = g.id[1]
	flake[11] = g.id[2]
	flake[12] = g.id[3]
	flake[13] = g.id[4]
	flake[14] = g.id[5]
	flake[15] = g.id[6]
	return flake
}

// ID returns the ID bytes associated with this generator.
func (g *Generator) ID() []byte {
	return g.id
}

// SID returns the secondary ID associated with this generator.
func (g *Generator) SID() int32 {
	return g.sid
}

// Time returns the Flake's timestamp as an int64.  The timestamp has
// microsecond resolution.
func (f *Flake) Time() int64 {
	return int64(f[0])<<33 | int64(f[1])<<25 | int64(f[2])<<17 | int64(f[3])<<9 | int64(f[4])<<1 + epoch
}

// ID returns the Flake's ID as a []byte.  A snowflake's ID is the last 8
// bytes.  No assumptios are made about either the contents of those bytes
// or their layout, that is left up to the user.
func (f *Flake) ID() []byte {
	return f[9:]
}

// SID returns the Flake's secondary ID.  The secondary ID is 20 bits.
func (f *Flake) SID() int32 {
	return int32(f[6]<<4>>4)<<16 | int32(f[7])<<8 | int32(f[8])
}

func seed() uint64 {
	bi := big.NewInt(int64(sequenceMax) + 1)
	r, err := rand.Int(rand.Reader, bi)
	if err != nil {
		panic(fmt.Sprintf("entropy read error: %s", err))
	}
	return r.Uint64()
}
