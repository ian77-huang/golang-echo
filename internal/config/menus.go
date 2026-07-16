package config

func SetMenus(rule MenuRules) []Menus {
	menus := []Menus{
		{Name: rule.T("index.title"), Url: "/"},
		{Name: rule.T("user.title"), Url: "/user", Childs: rule.users},
	}
	return menus
}
func SetMenusUsers(rule MenuUsersRules) []MenusChilds {
	users := []MenusChilds{}

	switch rule.Path {
	case "/user/login":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("user.register"), Url: "/user/register"})
		}
	case "/user/register":
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("user.login"), Url: "/user/login"})
		}
	case "/user/profile":
		if rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("user.reset.password.title"), Url: "/user/reset-password"})
		}
	case "/user/reset-password":
		if rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("user.profile.title"), Url: "/user/profile"})
		}
	default:
		if !rule.IsSignedIn {
			users = append(users, MenusChilds{Name: rule.T("user.login"), Url: "/user/login"})
			users = append(users, MenusChilds{Name: rule.T("user.register"), Url: "/user/register"})
		} else {
			users = append(users, MenusChilds{Name: rule.T("user.profile.title"), Url: "/user/profile"})
			users = append(users, MenusChilds{Name: rule.T("user.reset.password.title"), Url: "/user/reset-password"})
		}
	}
	return users
}
