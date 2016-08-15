package snoflinga

import (
	"bytes"
	"testing"
)

var tests = []struct {
	microseconds int64
	sequence     int32
	expected     Flake
}{
	{1257894000000496, 3499, [16]byte{0x47, 0x80, 0xc4, 0x50, 0x8f, 0xdf, 0x0d, 0xab, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1471274591071231, 4034, [16]byte{0x53, 0xa1, 0xdc, 0xf5, 0xe2, 0xbf, 0xff, 0xc2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1412312342545662, 0, [16]byte{0x50, 0x47, 0xd9, 0x77, 0xd4, 0x4f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1463234823749274, 124, [16]byte{0x53, 0x2c, 0xde, 0x7e, 0x47, 0xe9, 0xa0, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1321380928301928, 3212, [16]byte{0x4b, 0x1c, 0x9f, 0x8d, 0x82, 0xb6, 0x8c, 0x8c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	// 5
	{1382374987987923, 385, [16]byte{0x4e, 0x94, 0x34, 0x21, 0xaf, 0xbd, 0x31, 0x81, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1423898983276449, 42, [16]byte{0x50, 0xf0, 0x75, 0x11, 0x81, 0xfa, 0x10, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1449850923127583, 420, [16]byte{0x52, 0x6a, 0x1b, 0x94, 0x01, 0x31, 0xf1, 0xa4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1458375820399200, 1999, [16]byte{0x52, 0xe6, 0x29, 0x4b, 0x95, 0x26, 0x07, 0xcf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1467800048983349, 1024, [16]byte{0x53, 0x6f, 0x4d, 0x48, 0x55, 0xd3, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

// this tests an implementation of the time and sequence processing code
// implemented in a function local to the tests.  This is to validate that
// the actual implementation works as expected.  It's not reliable to test
// the actual implementation as it obtains the time itself.  This isn't
// necessarily reliable either as one must ensure that the actual
// implementation does not diverge with the one being tested here.  There is
// probably a better way, but this is good enough, for now.
func TestTimeSeq(t *testing.T) {
	for i, test := range tests {
		f := timeSequence(test.microseconds, test.sequence)
		if bytes.Compare(f[:], test.expected[:]) != 0 {
			t.Errorf("%d: got %x; want %x", i, f, test.expected)
			continue
		}
		// extract the time from it
		v := f.Time()
		if v != test.microseconds {
			t.Errorf("%d: got %d, want %d", i, v, test.microseconds)
		}
	}
}

// This only handles the first 8 bytes of a Flake - which is the time and
// sequence number portion.
func timeSequence(now int64, v int32) Flake {
	var flake Flake

	flake[0] = byte(now >> 44)
	flake[1] = byte(now >> 36)
	flake[2] = byte(now >> 28)
	flake[3] = byte(now >> 20)
	flake[4] = byte(now >> 12)
	flake[5] = byte(now >> 4)
	flake[6] = byte(now<<4) | uint8(v>>8)
	flake[7] = byte(v << 8 >> 8)

	return flake
}
