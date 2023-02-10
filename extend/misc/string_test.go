package misc

import "testing"

func TestString(t *testing.T) {
	s := "abcdefg"
	t.Log(s[0], rune(s[0]), string(s[0]), s[:4])
	t.Log(string('a'), string(rune(97)))

	s = "世界"
	t.Log(len(s))
	for _, e := range s {
		t.Log(e, string(e))
	}
	r := []rune(s)
	t.Log(r, len(r))
}
