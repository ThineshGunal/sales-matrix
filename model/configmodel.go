package model

type DBConnection struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
}
type DBConfig struct {
	MaxConnection     int `toml:"maxconn"`
	MaxIdleConnection int `toml:"maxidleconn"`
	MaxOpenConnection int `toml:"maxopenconn"`
}

type Config struct {
	DB       DBConnection `toml:"dbconnection"`
	DBConfig DBConfig     `toml:"dbconfig"`
}
