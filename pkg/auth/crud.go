package auth

func (a *Auth[TUser, TSession]) CreateUser(account string, password string) (*User[TUser], error) {
	config, err := a.getConfig()
	if err != nil {
		return nil, err
	}
	resolver := config.Resolver

	if resolver.CreateUser == nil {
		return nil, NewError("error.auth.InvalidConfiguration", "user creation resolver is not configured")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := resolver.CreateUser(account, passwordHash)
	if err != nil || user == nil || user.ID == "" {
		return nil, NewError("error.auth.CreateUserError", "create user error")
	}

	return user, nil
}
