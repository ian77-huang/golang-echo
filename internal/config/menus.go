package config

func SetMenus(rule MenuRules) []Menus {
	menus := []Menus{
		{Name: rule.T("index.title"), Url: "/"},
		{Name: rule.T("users.title"), Url: "/users", Childs: rule.users},
	}
	return menus
}
func SetMenusUsers(rule MenuUsersRules) []MenusChilds {
	users := []MenusChilds{}

	switch rule.Path {
	case "/users/login":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.register"), Url: "/users/register"})
		}
	case "/users/register":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.login"), Url: "/users/login"})
		}
	default:
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.register"), Url: "/users/register"})
			users = append(users, MenusChilds{Name: rule.T("users.login"), Url: "/users/login"})
		}
	}
	return users
}
