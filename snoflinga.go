// Package snoflinga generates snowflake like 128bit ids.  The first 52 bits
// is a timestamp representing time since Unix epoch, in microseconds.  The
// next 12 bits is a sequence number, that is increased with each snowflake
// request, for collision avoidance.  The start of the sequence is randomly
// selected.  The final 64 bits is the id.
package snoflinga

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	pcg "github.com/dgryski/go-pcgr"
)

const (
	idBits             = 64
	sequenceBits       = 12
	sequenceMax  int32 = 1 << sequenceBits
	// the sequenceBits only applies to the bits that are > 8; e.g for a 12bit
	// step the mask would apply to the first 8 bits.  No need to mask the last
	// 8 bits as they can be used as is.
	sequenceMask uint8 = -1 ^ (-1 << (sequenceBits % 8))
)

var (
	rng   pcg.Rand
	rngMu sync.Mutex
)

func init() {
	bi := big.NewInt(1<<63 - 1)
	r, err := rand.Int(rand.Reader, bi)
	if err != nil {
		panic(fmt.Sprintf("entropy read error: %s", err))
	}
	rng = pcg.New(r.Int64(), 0)
}

// ID is the type for a 128 bit snowflake.
type ID [16]byte

// A generator creates snowflakes for a given id
type Generator struct {
	id       []byte
	sequence int32
}

// New Generator returns an initialized generator.  If the passed byte slice is
// greater than 8 bytes, the first 8 bytes will be used for the generator's id.
// If the passed byte slice is less than 8 bytes, the id will be left-padded
// with 0, zero.  The generator's sequence is initialized with a random
// number.
func NewGenerator(id []byte) Generator {
	var g Generator
	if len(id) < 8 {
		g.id = make([]byte, 8-len(id))
	}
	g.id = append(g.id, id...)
	g.sequence = int32(rng.Bound(1<<sequenceBits - 1))
	return g
}

func (g *Generator) ID() ID {
	var flake ID
	now := uint64(time.Now().UnixNano() / 1000)
	v := atomic.AddInt32(&g.sequence, 1)
	if v == sequenceMax {
		v = atomic.AddInt32(&g.sequence, -v)
	}
	flake[0] = byte(now >> 44)
	flake[1] = byte(now >> 36)
	flake[2] = byte(now >> 28)
	flake[3] = byte(now >> 20)
	flake[4] = byte(now >> 12)
	flake[5] = byte(now >> 4)
	flake[6] = byte(now>>4<<4) | uint8(v>>8) ^ sequenceMask
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
