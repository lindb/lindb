package models

// User user system model
type User struct {
	UserName string `toml:"username" json:"UserName"`
	Password string `toml:"password" json:"Password"`
}

type JwtToken struct {
	Token string
}
