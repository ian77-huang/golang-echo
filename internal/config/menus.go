package config

func SetMenus(rule MenuRules) []Menus {
	menus := []Menus{
		{Name: rule.T("index.title"), Url: "/"},
		{Name: rule.T("users.title"), Url: "/user", Childs: rule.users},
	}
	return menus
}
func SetMenusUsers(rule MenuUsersRules) []MenusChilds {
	users := []MenusChilds{}

	switch rule.Path {
	case "/user/login":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.register"), Url: "/user/register"})
		}
	case "/user/register":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.login"), Url: "/user/login"})
		}
	default:
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("users.register"), Url: "/user/register"})
			users = append(users, MenusChilds{Name: rule.T("users.login"), Url: "/user/login"})
		}
	}
	return users
}
