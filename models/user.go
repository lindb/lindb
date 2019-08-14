package models

// User represents user model
type User struct {
	UserName string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}
