package math

import "testing"

func TestRandNumberFromStringDeterministic(t *testing.T) {
	for _, tt := range []struct {
		input string
		max   int
	}{{"hello", 10}, {"world", 100}, {"", 5}, {"中文", 20}, {"a-b-c", 1}} {
		first := RandNumberFromString(tt.input, tt.max)
		second := RandNumberFromString(tt.input, tt.max)
		if first != second {
			t.Fatalf("RandNumberFromString(%q, %d) not deterministic: %d != %d", tt.input, tt.max, first, second)
		}
		if first < 1 || first > tt.max {
			t.Fatalf("RandNumberFromString(%q, %d) = %d out of range [1, %d]", tt.input, tt.max, first, tt.max)
		}
	}
}

func TestRandNumberFromStringDifferentInputs(t *testing.T) {
	if RandNumberFromString("apple", 1000) == RandNumberFromString("banana", 1000) {
		t.Fatal("expected different inputs to produce different results")
	}
}

func TestCryptoNumberFromStringSHA(t *testing.T) {
	for _, tt := range []struct {
		input string
		max   int
	}{{"hello", 10}, {"world", 100}, {"", 5}, {"中文", 20}, {"a-b-c", 1}} {
		first := CryptoNumberFromStringSHA(tt.input, tt.max)
		second := CryptoNumberFromStringSHA(tt.input, tt.max)
		if first != second {
			t.Fatalf("CryptoNumberFromStringSHA(%q, %d) not deterministic: %d != %d", tt.input, tt.max, first, second)
		}
		if first < 1 || first > tt.max {
			t.Fatalf("CryptoNumberFromStringSHA(%q, %d) = %d out of range [1, %d]", tt.input, tt.max, first, tt.max)
		}
	}
}
