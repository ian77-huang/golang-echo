package validator

import "testing"

func TestValidateAndMessages(t *testing.T) {
	v := New()
	err := v.Validate(struct {
		Name string `validate:"required"`
		Kind string `validate:"oneof=cat dog"`
	}{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	messages := v.Messages(err)
	if messages["Name"] != "Name is required" || messages["Kind"] != "Kind must be one of: cat dog" {
		t.Fatalf("unexpected messages: %#v", messages)
	}
}

func TestMessagesForOrdinaryError(t *testing.T) {
	if New().Messages(assertError("bad"))["_"] != "bad" {
		t.Fatal("unexpected fallback")
	}
}

func TestMessagesUsesDefaultForOtherTags(t *testing.T) {
	err := New().Validate(struct {
		Age int `validate:"min=1"`
	}{})
	if err == nil || New().Messages(err)["Age"] != "Age is invalid" {
		t.Fatalf("unexpected messages: %#v", New().Messages(err))
	}
}

type assertError string

func (e assertError) Error() string { return string(e) }
