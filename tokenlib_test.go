package tokenlib

import "testing"

var (
	s = "I_LIKE_UNICORNS"
	d = map[string]interface{}{
		"data": map[string]interface{}{
			"sess":   "string",
			"userID": "42",
		},
	}
	e = 300
	t = "eyJkYXRhIjp7InNlc3MiOiJzdHJpbmciLCJ1c2VySUQiOiI0MiJ9LCJleHBpcmVzIjoxNTQ0NTgxMjk1LCJzYWx0IjoiT0dJM1lUWmwifYUhz0863Z66i3haGjeK3yzwHKkh6aBu5poMjHfKNbRY"
)

func TestMakeToken(t *testing.T) {
	_, e := MakeToken(d, s, int64(e))
	if e != nil {
		t.Errorf("Make token failed: %s", e.Error())
	}
}

func TestParseToken(t *testing.T) {
	tk, e := MakeToken(d, s, int64(e))
	if e != nil {
		t.Errorf("Make token failed: %s", e.Error())
	}
	_, e = ParseToken(tk, s)
	if e != nil {
		t.Errorf("Parse token failed: %s", e.Error())
	}
}

func BenchmarkMakeToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeToken(d, s, int64(e))
	}
}

func BenchmarkParseToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseToken(t, s)
	}
}
