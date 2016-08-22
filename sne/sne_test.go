package sne

import (
	"bytes"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		id          []byte
		sid         int32
		expectedID  []byte
		expectedSID int32
	}{
		{nil, 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{nil, 1, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1},
		{[]byte{}, 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{[]byte{}, 11, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 11},
		{[]byte(""), 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},

		{[]byte(""), 421, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 421},
		{[]byte("hello"), 0, []byte{0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, 0},
		{[]byte("hello"), 1127, []byte{0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, 1127},
		{[]byte("hello12"), 0, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 0},
		{[]byte("hello12"), 202, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 202},

		{[]byte("hello12"), 987403, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 987403},
		{[]byte("hello12abc"), 0, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 0},
		{[]byte("hello12abc"), 25500, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 25500},
		{[]byte("hello12abc"), 1492330, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32}, 443754}, //0x16 0xc5 0x6a > 0x06 0xc5 0x6a
	}
	for i, test := range tests {
		g := New(test.id, test.sid)
		if bytes.Compare(g.id, test.expectedID) != 0 {
			t.Errorf("%d: got %vl want %v", i, g.id, test.expectedID)
		}
		if g.sid != test.expectedSID {
			t.Errorf("%d: got %d; want %d", i, g.sid, test.expectedSID)
		}
	}
}

func TestFlake(t *testing.T) {
	id := "hello12"
	sid := int32(1211)
	g := New([]byte(id), sid)
	for i := 0; i < 10; i++ {
		// will just check against the second; assume that the milliseconds portion
		// is ok
		now := time.Now().Unix()
		f := g.Snowflake()
		if f.Time()/1000 != now {
			// if not equal, test against the next second; just in case a second boundary was passed
			if f.Time()/1000 != now+1 {
				t.Errorf("got %d; want %d", f.Time()/1000, now)
			}
			continue
		}
		if string(f.ID()) != id {
			t.Errorf("got %s, want %s", string(f.ID()), id)
			t.Errorf("got %x, want %x", string(f.ID()), id)
		}
		if f.SID() != sid {
			t.Errorf("got %d, want %d", f.SID(), sid)
		}
	}
}

func BenchmarkSnowFlake(b *testing.B) {
	g := New([]byte("test"), 42)
	for i := 0; i < b.N; i++ {
		g.Snowflake()
	}
}
