package cast

import "testing"

func TestStringToInt(t *testing.T) {
	for _, tt := range []struct {
		input     string
		def, want int
		wantErr   bool
	}{{"42", 1, 42, false}, {"", 7, 7, false}, {"bad", 7, 7, true}} {
		got, err := StringToInt(tt.input, tt.def)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Fatalf("StringToInt(%q) = %d, %v", tt.input, got, err)
		}
	}
}

func TestIntToString(t *testing.T) {
	if got := IntToString(-12); got != "-12" {
		t.Fatalf("got %q", got)
	}
}
