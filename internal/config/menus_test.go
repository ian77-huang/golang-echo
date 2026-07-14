package config

import "testing"

func TestSetMenusUsers(t *testing.T) {
	tf := func(id string, _ ...any) string { return id }
	for _, tt := range []struct {
		path     string
		signedIn bool
		want     int
	}{{"/user/login", false, 1}, {"/user/register", false, 1}, {"/", false, 2}, {"/", true, 0}} {
		got := SetMenusUsers(MenuUsersRules{Path: tt.path, IsSignedIn: tt.signedIn, T: tf})
		if len(got) != tt.want {
			t.Fatalf("%s signedIn=%v: got %#v", tt.path, tt.signedIn, got)
		}
	}
}

func TestSetMenusIncludesUserMenu(t *testing.T) {
	menus := SetMenus(MenuRules{users: []MenusChilds{{Name: "login", Url: "/user/login"}}, T: func(id string, _ ...any) string { return id }})
	if len(menus) != 2 || menus[1].Childs[0].Name != "login" {
		t.Fatalf("unexpected menus: %#v", menus)
	}
}
