package model

import "testing"

func TestTableNames(t *testing.T) {
	for _, tt := range []struct {
		name string
		got  string
		want string
	}{
		{"User", (User{}).TableName(), "user"},
		{"UserProfile", (UserProfile{}).TableName(), "user_profile"},
		{"Session", (Session{}).TableName(), "session"},
		{"Bible", (Bible{}).TableName(), "bible"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("TableName() = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
