package argon2

import (
	"errors"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		password string
		want     bool
	}{{"correct horse", true}, {"wrong", false}} {
		got, err := VerifyPassword(tt.password, hash)
		if err != nil || got != tt.want {
			t.Fatalf("VerifyPassword() = %v, %v", got, err)
		}
	}
}

func TestVerifyPasswordRejectsMalformedHash(t *testing.T) {
	if _, err := VerifyPassword("x", "invalid"); err == nil {
		t.Fatal("expected error")
	}
}

func TestVerifyPasswordRejectsInvalidEncodedParts(t *testing.T) {
	for _, encoded := range []string{
		"$argon2id$v=19$m=bad,t=3,p=2$salt$hash",
		"$argon2id$v=19$m=1,t=1,p=1$%%%$hash",
		"$argon2id$v=19$m=1,t=1,p=1$c2FsdA$%%%",
	} {
		if _, err := VerifyPassword("x", encoded); err == nil {
			t.Fatalf("expected error for %q", encoded)
		}
	}
}

func TestHashPasswordReturnsRandomSourceError(t *testing.T) {
	original := RandomRead
	RandomRead = func([]byte) (int, error) { return 0, errors.New("random failed") }
	t.Cleanup(func() { RandomRead = original })
	if _, err := HashPassword("password"); err == nil {
		t.Fatal("expected random source error")
	}
}
