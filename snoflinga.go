// Package snöflinga generates snowflake like 128-bit ids.  The first 52 bits
// is a timestamp representing time since Unix epoch, in microseconds.  The
// next 12 bits is a sequence number which is increased with each snowflake
// request, for collision avoidance.  The start of the sequence is randomly
// selected.  The final 64 bits is the id.
package snoflinga

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"
)

const (
	idBits              = 64
	sequenceBits        = 12
	sequenceMax  uint64 = 1<<sequenceBits - 1
)

// Flake is the type for a 128 bit snowflake.
type Flake [16]byte

// Generator creates snowflakes for a given id
type Generator struct {
	id       []byte
	sequence uint64
}

// New returns an initialized generator.  If the passed byte slice's length is
// greater than 8 bytes, the first 8 bytes will be used for the generator's id.
// If the passed byte slice's is less than 8 bytes, the id will be left-padded
// with 0s, zeros.  The generator's sequence is initialized with a random
// number.
func New(id []byte) Generator {
	var g Generator
	if len(id) < 8 {
		g.id = make([]byte, 8-len(id))
	}
	if len(id) > 8 {
		id = id[:8]
	}
	g.id = append(g.id, id...)
	g.sequence = seed()
	return g
}

// Snowflake generates an Flake from the current time and next sequence value.
func (g *Generator) Snowflake() Flake {
	var flake Flake
	now := uint64(time.Now().UnixNano() / 1000)
	v := atomic.AddUint64(&g.sequence, 1)
	v = v % sequenceMax
	flake[0] = byte(now >> 44)
	flake[1] = byte(now >> 36)
	flake[2] = byte(now >> 28)
	flake[3] = byte(now >> 20)
	flake[4] = byte(now >> 12)
	flake[5] = byte(now >> 4)
	flake[6] = byte(now<<4) | uint8(v>>8)
	flake[7] = byte(v << 8 >> 8)
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

// Time returns the Flake's timestamp, in microseconds, as an int64.
func (f *Flake) Time() int64 {
	return int64(f[0])<<44 | int64(f[1])<<36 | int64(f[2])<<28 | int64(f[3])<<20 | int64(f[4])<<12 | int64(f[5])<<4 | int64(f[6]>>4)
}

// ID returns the Flake's ID as a []byte.  A snowflake's ID is the last 8
// bytes of the Flake.  No assumptiosn are made about either the contents of
// those bytes or their layout: that is left up to the user.
func (f *Flake) ID() []byte {
	return f[8:]
}

func seed() uint64 {
	bi := big.NewInt(int64(sequenceMax) + 1)
	r, err := rand.Int(rand.Reader, bi)
	if err != nil {
		panic(fmt.Sprintf("entropy read error: %s", err))
	}
	return r.Uint64()
}
