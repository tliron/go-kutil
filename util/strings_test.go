package util

import (
	"bytes"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	s := "this is a string"
	b1 := StringToBytes(s)
	b2 := ([]byte)(s)
	if !bytes.Equal(b1, b2) {
		t.Error("StringToBytes")
	}
	if StringToBytes("") != nil {
		t.Error("StringToBytes nil")
	}
}

func TestBytesToString(t *testing.T) {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	s1 := BytesToString(b)
	s2 := (string)(b)
	if s1 != s2 {
		t.Error("BytesToString")
	}
	if BytesToString(nil) != "" {
		t.Error("BytesToString nil")
	}
}

func TestJoinQuote(t *testing.T) {
	couple := []string{"hello", `"world"`}
	many := []string{"hello", `"world"`, "from", "Tal"}
	if JoinQuote(couple, ", ") != `"hello", "\"world\""` {
		t.Error("JoinQuote couple")
	}
	if JoinQuote(many, ", ") != `"hello", "\"world\"", "from", "Tal"` {
		t.Error("JoinQuote many")
	}
	if JoinQuoteL(couple, ", ", ", and ", " and ") != `"hello" and "\"world\""` {
		t.Error("JoinQuoteL couple")
	}
	if JoinQuoteL(many, ", ", ", and ", " and ") != `"hello", "\"world\"", "from", and "Tal"` {
		t.Error("JoinQuoteL many")
	}
}
