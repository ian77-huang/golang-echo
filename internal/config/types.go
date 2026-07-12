package config

type Menus struct {
	Name   string
	Url    string
	Childs []MenusChilds
}
type MenusChilds struct {
	Name string
	Url  string
}

type ConfigDatabases struct {
	Path string
}
type ConfigUsers struct {
	MinLengthAccount  int
	MinLengthPassword int
}

type Config struct {
	SecretKey string
	Databases ConfigDatabases
	Users     ConfigUsers
}
