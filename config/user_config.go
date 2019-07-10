package config

type LinDbUser struct {
	USER USER `toml:"USER"`
}

type USER struct {
	UserName string `toml:"UserName"`
	Password string `toml:"Password"`
}
