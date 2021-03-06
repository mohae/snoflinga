package sno

import (
	"bytes"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		id          []byte
		sid         int16
		expectedID  []byte
		expectedSID int16
	}{
		{nil, 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{nil, 1, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1},
		{[]byte{}, 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{[]byte{}, 11, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 11},
		{[]byte(""), 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},

		{[]byte(""), 42, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 42},
		{[]byte("hello"), 0, []byte{0x00, 0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, 0},
		{[]byte("hello"), 127, []byte{0x00, 0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, 127},
		{[]byte("hello123"), 0, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 0},
		{[]byte("hello123"), 202, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 202},

		{[]byte("hello123"), 9403, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 1211}, //0x24 0xbb > 0x04 0xbb
		{[]byte("hello123abc"), 0, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 0},
		{[]byte("hello123abc"), 255, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 255},
		{[]byte("hello123abc"), 14923, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x31, 0x32, 0x33}, 6731}, //0x3a 0x4b > 0x1a 0x4b
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
	id := "hello123"
	sid := int16(1211)
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
