package models

// User represents user model
type User struct {
	UserName string `toml:"username" json:"UserName"`
	Password string `toml:"password" json:"Password"`
}
